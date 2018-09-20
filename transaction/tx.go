package transaction

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"sort"
	"strconv"

	// "crypto/sha256"
	"math/big"

	"time"

	"unsafe"

	"github.com/ninjadotorg/cash-prototype/cashec"
	"github.com/ninjadotorg/cash-prototype/common"
	"github.com/ninjadotorg/cash-prototype/privacy/client"
	"github.com/ninjadotorg/cash-prototype/privacy/proto/zksnark"
)

// Tx represents a coin-transfer-transaction stored in a block
type Tx struct {
	Version  int8   `json:"Version"`
	Type     string `json:"Type"` // n
	LockTime int64  `json:"LockTime"`
	Fee      uint64 `json:"Fee"`

	Descs    []*JoinSplitDesc `json:"Descs"`
	JSPubKey []byte           `json:"JSPubKey,omitempty"` // 64 bytes
	JSSig    []byte           `json:"JSSig,omitempty"`    // 64 bytes

	txId       *common.Hash
	sigPrivKey *client.PrivateKey
}

func (tx *Tx) SetTxId(txId *common.Hash) {
	tx.txId = txId
}

func (tx *Tx) GetTxId() *common.Hash {
	return tx.txId
}

// Hash returns the hash of all fields of the transaction
func (tx Tx) Hash() *common.Hash {
	record := strconv.Itoa(int(tx.Version))
	record += tx.Type
	record += strconv.Itoa(int(tx.LockTime))
	record += strconv.Itoa(len(tx.Descs))
	for _, desc := range tx.Descs {
		record += desc.toString()
	}
	record += string(tx.JSPubKey)
	record += string(tx.JSSig)
	hash := common.DoubleHashH([]byte(record))
	return &hash
}

// ValidateTransaction returns true if transaction is valid:
// - JSDescriptions are valid (zk-snark proof satisfied)
// - Signature matches the signing public key
// Note: This method doesn't check for double spending
func (tx *Tx) ValidateTransaction() bool {
	for _, desc := range tx.Descs {
		if desc.Reward != 0 {
			return false // Coinbase tx shouldn't be broadcasted across the network
		}
	}

	// TODO: implement
	return true
}

// GetType returns the type of the transaction
func (tx *Tx) GetType() string {
	return tx.Type
}

// GetTxVirtualSize computes the virtual size of a given transaction
func (tx *Tx) GetTxVirtualSize() uint64 {
	return uint64(unsafe.Sizeof(tx))
}

// NewTxTemplate returns a new Tx initialized with default data
func NewTxTemplate() *Tx {
	//Generate signing key 96 bytes
	sigPrivKey, err := client.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	// Verification key 64 bytes
	sigPubKey := PubKeyToByteArray(&sigPrivKey.PublicKey)

	tx := &Tx{
		Version:  TxVersion,
		Type:     common.TxNormalType,
		LockTime: time.Now().Unix(),
		Fee:      0,
		Descs:    nil,
		JSPubKey: sigPubKey,
		JSSig:    nil,

		txId:       nil,
		sigPrivKey: sigPrivKey,
	}
	return tx
}

// CreateTx creates transaction with appropriate proof for a private payment
// value: total value of the coins to transfer
// rt: root of the commitment merkle tree at current block (the latest block of the node creating this tx)
func CreateTx(
	senderKey *client.SpendingKey,
	paymentInfo []*client.PaymentInfo,
	rt *common.Hash,
	usableTx []*Tx,
	nullifiers [][]byte,
	commitments [][]byte,
) (*Tx, error) {
	fmt.Printf("List of all commitments before building tx:\n")
	for _, cm := range commitments {
		fmt.Printf("%x\n", cm)
	}

	var value uint64
	for _, p := range paymentInfo {
		value += p.Amount
		fmt.Printf("[CreateTx] paymentInfo.Value: %v, paymentInfo.Apk: %x\n", p.Amount, p.PaymentAddress.Apk)
	}

	// Get list of notes to use
	var inputNotes []*client.Note
	for _, tx := range usableTx {
		for _, desc := range tx.Descs {
			for _, note := range desc.note {
				inputNotes = append(inputNotes, note)
				fmt.Printf("[CreateTx] inputNote.Value: %v\n", note.Value)
			}
		}
	}

	// Left side value
	var sumInputValue uint64
	for _, note := range inputNotes {
		sumInputValue += note.Value
	}
	if sumInputValue < value {
		return nil, fmt.Errorf("Input value less than output value")
	}

	senderFullKey := cashec.KeySet{}
	senderFullKey.ImportFromPrivateKeyByte(senderKey[:])

	// Create tx before adding js descs
	tx := NewTxTemplate()
	latestAnchor := []byte{}

	for len(inputNotes) > 0 || len(paymentInfo) > 0 {
		// Sort input and output notes ascending by value to start building js descs
		sort.Slice(inputNotes, func(i, j int) bool {
			return inputNotes[i].Value < inputNotes[j].Value
		})
		sort.Slice(paymentInfo, func(i, j int) bool {
			return paymentInfo[i].Amount < paymentInfo[j].Amount
		})

		// Choose inputs to build js desc
		var inputsToBuildWitness, inputs []*client.JSInput
		var inputValue uint64
		for len(inputNotes) > 0 && len(inputs) < NumDescInputs {
			input := &client.JSInput{}
			input.InputNote = inputNotes[len(inputNotes)-1] // Get note with largest value
			input.Key = senderKey
			inputs = append(inputs, input)
			inputsToBuildWitness = append(inputsToBuildWitness, input)
			inputValue += input.InputNote.Value

			inputNotes = inputNotes[:len(inputNotes)-1]
			fmt.Printf("Choose input note with value %v and cm %x\n", input.InputNote.Value, input.InputNote.Cm)
		}

		// Add dummy input note if necessary
		for len(inputs) < NumDescInputs {
			input := &client.JSInput{}
			input.InputNote = createDummyNote(senderKey)
			input.Key = senderKey
			input.WitnessPath = (&client.MerklePath{}).CreateDummyPath() // No need to build commitment merkle path for dummy note
			inputs = append(inputs, input)
			fmt.Printf("Add dummy input note\n")
		}

		// Check if input note's cm is in commitments list
		for _, input := range inputsToBuildWitness {
			input.InputNote.Cm = client.GetCommitment(input.InputNote)

			found := false
			for _, c := range commitments {
				if bytes.Equal(c, input.InputNote.Cm) {
					found = true
				}
			}
			if found == false {
				return nil, fmt.Errorf("Commitment %x of input note isn't in commitments list", input.InputNote.Cm)
			}
		}

		// Build witness path for the input notes
		newRt, err := client.BuildWitnessPath(inputsToBuildWitness, commitments)
		if err != nil {
			return nil, err
		}

		// For first js desc, check if provided Rt is the root of the merkle tree contains all commitments
		if latestAnchor == nil {
			if !bytes.Equal(newRt, rt[:]) {
				return nil, fmt.Errorf("Provided anchor doesn't match commitments list")
			}
		}
		latestAnchor = newRt

		// Choose output notes for the js desc
		outputs := []*client.JSOutput{}
		for len(paymentInfo) > 0 && len(outputs) < NumDescOutputs-1 && inputValue > 0 { // Leave out 1 output note for change
			p := paymentInfo[len(paymentInfo)-1]
			var outNote *client.Note
			var encKey client.TransmissionKey
			if p.Amount <= inputValue { // Enough for one more output note, include it
				outNote = &client.Note{Value: p.Amount, Apk: p.PaymentAddress.Apk}
				encKey = p.PaymentAddress.Pkenc
				inputValue -= p.Amount
				paymentInfo = paymentInfo[:len(paymentInfo)-1]
				fmt.Printf("Use output value %v => %x\n", outNote.Value, outNote.Apk)
			} else { // Not enough for this note, send some and save the rest for next js desc
				outNote = &client.Note{Value: inputValue, Apk: p.PaymentAddress.Apk}
				encKey = p.PaymentAddress.Pkenc
				paymentInfo[len(paymentInfo)-1].Amount = p.Amount - inputValue
				inputValue = 0
				fmt.Printf("Partially send %v to %x\n", outNote.Value, outNote.Apk)
			}

			output := &client.JSOutput{EncKey: encKey, OutputNote: outNote}
			outputs = append(outputs, output)
		}

		if inputValue > 0 {
			// Still has some room left, check if one more output note is possible to add
			var p *client.PaymentInfo
			if len(paymentInfo) > 0 {
				p = paymentInfo[len(paymentInfo)-1]
			}

			if p != nil && p.Amount == inputValue {
				// Exactly equal, add this output note to js desc
				outNote := &client.Note{Value: p.Amount, Apk: p.PaymentAddress.Apk}
				output := &client.JSOutput{EncKey: p.PaymentAddress.Pkenc, OutputNote: outNote}
				outputs = append(outputs, output)
				paymentInfo = paymentInfo[:len(paymentInfo)-1]
				fmt.Printf("Exactly enough, include 1 more output %v, %x\n", outNote.Value, outNote.Apk)
			} else {
				// Cannot put the output note into this js desc, create a change note instead
				outNote := &client.Note{Value: inputValue, Apk: senderFullKey.PublicKey.Apk}
				output := &client.JSOutput{EncKey: senderFullKey.PublicKey.Pkenc, OutputNote: outNote}
				outputs = append(outputs, output)
				fmt.Printf("Create change outnote %v, %x\n", outNote.Value, outNote.Apk)

				// Use the change note to continually send to receivers if needed
				if len(paymentInfo) > 0 {
					// outNote data (R and Rho) will be updated when building zk-proof
					inputNotes = append(inputNotes, outNote)
					fmt.Printf("Reuse change note later\n")
				}
			}
			inputValue = 0
		}

		// Add dummy output note if necessary
		for len(outputs) < NumDescOutputs {
			outputs = append(outputs, CreateRandomJSOutput())
			fmt.Printf("Create dummy output note\n")
		}

		// TODO: Shuffle output notes randomly (if necessary)

		// Generate proof and sign tx
		var reward uint64 // Zero reward for non-coinbase transaction
		err = tx.BuildNewJSDesc(inputs, outputs, latestAnchor, reward)
		if err != nil {
			return nil, err
		}

		// Add new commitments to list to use in next js desc if needed
		for _, output := range outputs {
			fmt.Printf("Add new output cm to list: %x\n", output.OutputNote.Cm)
			commitments = append(commitments, output.OutputNote.Cm)
		}

		fmt.Printf("Len input and info: %v %v\n", len(inputNotes), len(paymentInfo))
	}

	// Sign tx
	tx, err := SignTx(tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("jspubkey size: %v\n", len(tx.JSPubKey))
	fmt.Printf("jssig size: %v\n", len(tx.JSSig))
	return tx, nil
}

// BuildNewJSDesc creates zk-proof for a js desc and add it to the transaction
func (tx *Tx) BuildNewJSDesc(inputs []*client.JSInput, outputs []*client.JSOutput, rt []byte, reward uint64) error {
	var seed, phi []byte
	var outputR [][]byte
	proof, hSig, seed, phi, err := client.Prove(inputs, outputs, tx.JSPubKey, rt, reward, seed, phi, outputR)
	if err != nil {
		return err
	}

	var ephemeralPrivKey *client.EphemeralPrivKey // nil ephemeral key, will be randomly created later
	err = tx.buildJSDescAndEncrypt(inputs, outputs, proof, rt, reward, hSig, seed, tx.JSPubKey, ephemeralPrivKey)
	if err != nil {
		return err
	}
	return nil
}

func (tx *Tx) buildJSDescAndEncrypt(
	inputs []*client.JSInput,
	outputs []*client.JSOutput,
	proof *zksnark.PHGRProof,
	rt []byte,
	reward uint64,
	hSig, seed, sigPubKey []byte,
	ephemeralPrivKey *client.EphemeralPrivKey,
) error {
	nullifiers := [][]byte{inputs[0].InputNote.Nf, inputs[1].InputNote.Nf}
	commitments := [][]byte{outputs[0].OutputNote.Cm, outputs[1].OutputNote.Cm}
	notes := [2]client.Note{*outputs[0].OutputNote, *outputs[1].OutputNote}
	keys := [2]client.TransmissionKey{outputs[0].EncKey, outputs[1].EncKey}

	ephemeralPubKey := new(client.EphemeralPubKey)
	if ephemeralPrivKey == nil {
		ephemeralPrivKey = new(client.EphemeralPrivKey)
		*ephemeralPubKey, *ephemeralPrivKey = client.GenEphemeralKey()
	} else { // Genesis block only
		ephemeralPrivKey.GenPubKey()
		*ephemeralPubKey = ephemeralPrivKey.GenPubKey()
	}
	fmt.Printf("hSig: %x\n", hSig)
	fmt.Printf("ephemeralPrivKey: %x\n", *ephemeralPrivKey)
	fmt.Printf("ephemeralPubKey: %x\n", *ephemeralPubKey)
	fmt.Printf("tranmissionKey[0]: %x\n", keys[0])
	fmt.Printf("tranmissionKey[1]: %x\n", keys[1])
	fmt.Printf("notes[0].Value: %v\n", notes[0].Value)
	fmt.Printf("notes[0].Rho: %x\n", notes[0].Rho)
	fmt.Printf("notes[0].R: %x\n", notes[0].R)
	fmt.Printf("notes[0].Memo: %v\n", notes[0].Memo)
	fmt.Printf("notes[1].Value: %v\n", notes[1].Value)
	fmt.Printf("notes[1].Rho: %x\n", notes[1].Rho)
	fmt.Printf("notes[1].R: %x\n", notes[1].R)
	fmt.Printf("notes[1].Memo: %v\n", notes[1].Memo)
	noteciphers := client.EncryptNote(notes, keys, *ephemeralPrivKey, *ephemeralPubKey, hSig)

	//Calculate vmacs to prove this transaction is signed by this user
	vmacs := make([][]byte, 2)
	for i := range inputs {
		ask := make([]byte, 32)
		copy(ask[:], inputs[i].Key[:])
		vmacs[i] = client.PRF_pk(uint64(i), ask, hSig)
	}

	desc := &JoinSplitDesc{
		Anchor:          rt,
		Nullifiers:      nullifiers,
		Commitments:     commitments,
		Proof:           proof,
		EncryptedData:   noteciphers,
		EphemeralPubKey: ephemeralPubKey[:],
		HSigSeed:        seed,
		Type:            common.TxOutCoinType,
		Reward:          reward,
		Vmacs:           vmacs,
	}
	tx.Descs = append(tx.Descs, desc)

	fmt.Println("desc:")
	fmt.Printf("Anchor: %x\n", desc.Anchor)
	fmt.Printf("Nullifiers: %x\n", desc.Nullifiers)
	fmt.Printf("Commitments: %x\n", desc.Commitments)
	fmt.Printf("Proof: %x\n", desc.Proof)
	fmt.Printf("EncryptedData: %x\n", desc.EncryptedData)
	fmt.Printf("EphemeralPubKey: %x\n", desc.EphemeralPubKey)
	fmt.Printf("HSigSeed: %x\n", desc.HSigSeed)
	fmt.Printf("Type: %v\n", desc.Type)
	fmt.Printf("Reward: %v\n", desc.Reward)
	fmt.Printf("Vmacs: %x %x\n", desc.Vmacs[0], desc.Vmacs[1])
	return nil
}

// CreateRandomJSInput creates a dummy input with 0 value note that belongs to a random address
func CreateRandomJSInput() *client.JSInput {
	randomKey := client.RandSpendingKey()
	input := new(client.JSInput)
	input.InputNote = createDummyNote(&randomKey)
	input.Key = &randomKey
	input.WitnessPath = (&client.MerklePath{}).CreateDummyPath()
	return input
}

// CreateRandomJSOutput creates a dummy output with 0 value note that is sended to a random address
func CreateRandomJSOutput() *client.JSOutput {
	randomKey := client.RandSpendingKey()
	output := new(client.JSOutput)
	output.OutputNote = createDummyNote(&randomKey)
	paymentAddr := client.GenPaymentAddress(randomKey)
	output.EncKey = paymentAddr.Pkenc
	return output
}

func createDummyNote(spendingKey *client.SpendingKey) *client.Note {
	addr := client.GenSpendingAddress(*spendingKey)
	var rho, r [32]byte
	copy(rho[:], client.RandBits(32*8))
	copy(r[:], client.RandBits(32*8))

	note := &client.Note{
		Value: 0,
		Apk:   addr,
		Rho:   rho[:],
		R:     r[:],
		Nf:    client.GetNullifier(*spendingKey, rho),
	}
	return note
}

func SignTx(tx *Tx) (*Tx, error) {
	//Check input transaction
	if tx.JSSig != nil {
		return nil, errors.New("Input transaction must be an unsigned one")
	}

	// Hash transaction
	tx.SetTxId(tx.Hash())
	hash := tx.GetTxId()
	data := make([]byte, common.HashSize)
	copy(data, hash[:])

	// Sign
	ecdsaSignature := new(client.EcdsaSignature)
	var err error
	ecdsaSignature.R, ecdsaSignature.S, err = client.Sign(rand.Reader, tx.sigPrivKey, data[:])
	if err != nil {
		return nil, err
	}

	//Signature 64 bytes
	tx.JSSig = JSSigToByteArray(ecdsaSignature)

	return tx, nil
}

func VerifySign(tx *Tx) (bool, error) {
	//Check input transaction
	if tx.JSSig == nil || tx.JSPubKey == nil {
		return false, errors.New("Input transaction must be an signed one!")
	}

	// UnParse Public key
	pubKey := new(client.PublicKey)
	pubKey.X = new(big.Int).SetBytes(tx.JSPubKey[0:32])
	pubKey.Y = new(big.Int).SetBytes(tx.JSPubKey[32:64])

	// UnParse ECDSA signature
	ecdsaSignature := new(client.EcdsaSignature)
	ecdsaSignature.R = new(big.Int).SetBytes(tx.JSSig[0:32])
	ecdsaSignature.S = new(big.Int).SetBytes(tx.JSSig[32:64])

	// Hash origin transaction
	hash := tx.GetTxId()
	data := make([]byte, common.HashSize)
	copy(data, hash[:])

	valid := client.VerifySign(pubKey, data[:], ecdsaSignature.R, ecdsaSignature.S)
	return valid, nil
}

// GenerateProofForGenesisTx creates zk-proof and build the transaction (without signing) for genesis block
func GenerateProofForGenesisTx(
	inputs []*client.JSInput,
	outputs []*client.JSOutput,
	rt []byte,
	reward uint64,
	seed, phi []byte,
	outputR [][]byte,
	ephemeralPrivKey client.EphemeralPrivKey,
) (*Tx, error) {
	// Generate JoinSplit key pair to act as a dummy key (since we don't sign genesis tx)
	privateSignKey := [32]byte{1}
	keyPair := &cashec.KeySet{}
	keyPair.ImportFromPrivateKeyByte(privateSignKey[:])
	sigPubKey := keyPair.PublicKey.Apk[:]

	tx := NewTxTemplate()
	tx.JSPubKey = sigPubKey

	proof, hSig, seed, phi, err := client.Prove(inputs, outputs, tx.JSPubKey, rt, reward, seed, phi, outputR)
	if err != nil {
		return nil, err
	}

	err = tx.buildJSDescAndEncrypt(inputs, outputs, proof, rt, reward, hSig, seed, tx.JSPubKey, &ephemeralPrivKey)
	return tx, err
}

func PubKeyToByteArray(pubKey *client.PublicKey) []byte {
	var pub []byte
	pubX := pubKey.X.Bytes()
	pubY := pubKey.Y.Bytes()
	pub = append(pub, pubX...)
	pub = append(pub, pubY...)
	return pub
}

func JSSigToByteArray(jsSig *client.EcdsaSignature) []byte {
	var jssig []byte
	r := jsSig.R.Bytes()
	s := jsSig.S.Bytes()
	jssig = append(jssig, r...)
	jssig = append(jssig, s...)
	return jssig
}

func SortArrayTxs(data []Tx, sortType int, sortAsc bool) {
	if len(data) == 0 {
		return
	}
	switch sortType {
	case NoSort:
		{
			// do nothing
		}
	case SortByAmount:
		{
			sort.SliceStable(data, func(i, j int) bool {
				desc1 := data[i].Descs
				amount1 := uint64(0)
				for _, desc := range desc1 {
					for _, note := range desc.GetNote() {
						amount1 += note.Value
					}
				}
				desc2 := data[j].Descs
				amount2 := uint64(0)
				for _, desc := range desc2 {
					for _, note := range desc.GetNote() {
						amount2 += note.Value
					}
				}
				if !sortAsc {
					return amount1 > amount2
				} else {
					return amount1 <= amount2
				}
			})
		}
	default:
		{
			// do nothing
		}
	}
}
