package crypto

// import (
// 	"os"

// 	"github.com/tidwall/go-node"
// )

// type CryptoService struct {
// 	vm node.VM

// 	password       string
// 	firstNonce     string
// 	lastNonce      string
// 	salt           string
// 	scarmSalt      string
// 	iterations     int
// 	saltedPassword string
// }

// func NewCryptoService() *CryptoService {
// 	password := os.Getenv("PASSWORD")
// 	if password == "" {
// 		panic("password not set")
// 	}

// 	nodeOps := &node.Options{}
// 	vm := node.New(nodeOps)
// 	vm.Run(nodePackage)

// 	newCryptoService := &CryptoService{
// 		vm: vm,

// 		password:   password,
// 		iterations: 100,
// 	}

// 	firstNonce, err := newCryptoService.getFirstNonce()
// 	if err != nil {
// 		panic(err)
// 	}

// 	newCryptoService.firstNonce = *firstNonce

// 	return newCryptoService
// }

// func (c *CryptoService) runCryptoCommand(command string) (any, error) {
// 	v := c.vm.Run(command)
// 	if v.Error() != nil {
// 		return nil, v.Error()
// 	}

// 	output := v.String()
// 	return &output, nil
// }

// func (c *CryptoService) getFirstNonce() (*string, error) {
// 	command := `scram.nonce().toString();`
// 	outObj, err := c.runCryptoCommand(command)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if outObj == nil {
// 		return nil, fmt.Errorf("first nonce is nil")
// 	}

// 	outStr := outObj.(*string)

// 	return outStr, nil
// }

// func (c *CryptoService) getScarmSalt() error {
// 	if c.salt == "" {
// 		log.Warn().Msg("salt not set")
// 		return nil
// 	}

// 	command := fmt.Sprintf(`cryptoJS.enc.Hex.parse("%s");`, c.salt)
// 	outObj, err := c.runCryptoCommand(command)
// 	if err != nil {
// 		return err
// 	}

// 	if outObj == nil {
// 		return fmt.Errorf("output is nil")
// 	}

// 	outStr := outObj.(*string)

// 	c.scarmSalt = *outStr

// 	return nil
// }

// func (c *CryptoService) getSaltedPassword() error {
// 	err := c.getScarmSalt()
// 	if err != nil {
// 		return err
// 	}

// 	command := fmt.Sprintf(`scram.saltedPassword("%s", "%s", "%d").toString();`, c.password, c.scarmSalt, c.iterations)
// 	outObj, err := c.runCryptoCommand(command)
// 	if err != nil {
// 		return err
// 	}

// 	if outObj == nil {
// 		return fmt.Errorf("output is nil")
// 	}

// 	outStr := outObj.(*string)
// 	c.saltedPassword = *outStr

// 	return nil
// }

// func (c *CryptoService) getAuthData() (*string, error) {
// 	authMessage := fmt.Sprintf("%s,%s,%s", c.firstNonce, c.lastNonce, c.lastNonce)

// 	command := fmt.Sprintf(`scram.clientProof("%s", "%s", "%d", "%s").toString();`, c.password, c.scarmSalt, c.iterations, authMessage)

// 	outObj, err := c.runCryptoCommand(command)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if outObj == nil {
// 		return nil, fmt.Errorf("out is nil")
// 	}

// 	outStr := outObj.(*string)

// 	return outStr, nil
// }
