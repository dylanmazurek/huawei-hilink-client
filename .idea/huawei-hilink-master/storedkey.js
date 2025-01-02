const huawei = require('/home/dylan/dev/huawei-wifi/.idea/huawei-hilink-master/jslib/public');
const CryptoJS = huawei.CryptoJS;

function testScram() {
    const password = "testpass";
    const salt = "4ec91c5b362986cec5ef6db10eb4bb5b1c54eef1eff89006cf00a8993226fabd";
    const iter = 100;

    const scram = CryptoJS.SCRAM();
    const scarmSalt = CryptoJS.enc.Hex.parse(salt);
    
    // Step 1: Salt Password
    const saltPassword = scram.saltedPassword(password, scarmSalt, iter);
    console.log('saltPassword:', saltPassword.toString(CryptoJS.enc.Hex));
    
    // Step 2: Client Key
    const clientKey = scram.clientKey(saltPassword);
    console.log('clientKey:', clientKey.toString(CryptoJS.enc.Hex));
    
    // Step 3: Stored Key (plain SHA256)
    const storedKey = scram.storedKey(clientKey);
    console.log('storedKey:', storedKey.toString(CryptoJS.enc.Hex));
}

testScram();