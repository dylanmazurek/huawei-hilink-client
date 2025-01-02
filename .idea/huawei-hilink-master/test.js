const huawei = require('/home/dylan/dev/huawei-wifi/.idea/huawei-hilink-master/jslib/public');

var password = 'testpass';
var salt = '4ec91c5b362986cec5ef6db10eb4bb5b1c54eef1eff89006cf00a8993226fabd';
var nonce = '868f304eec646840959cb3830c0fe352bbc3b24f264ddd652764c8ad2c5fabe3';
var finalNonce = '868f304eec646840959cb3830c0fe352bbc3b24f264ddd652764c8ad2c5fabe31RYXGq1AJapOUsCb57h0bT5oCTsgDIDu';

function main(){
    var cryptoJS = huawei.CryptoJS;
    var scram = cryptoJS.SCRAM();
    
    console.log(`salt: ${salt}`);

    const iter = 100;
    
    var saltPassword = scram.saltedPassword(password, scarmSalt, iter);
    console.log(`saltPassword: ${saltPassword.toString(cryptoJS.enc.Hex)}`);
    
    var clientKey = scram.clientKey(saltPassword).toString(cryptoJS.enc.Hex);
    console.log(`clientKey: ${clientKey.toString(cryptoJS.enc.Hex)}`);

    var serverKey = scram.serverKey(saltPassword).toString(cryptoJS.enc.Hex);
    console.log(`serverKey: ${serverKey.toString(cryptoJS.enc.Hex)}`);

    var storedKey = scram.storedKey(clientKey).toString(cryptoJS.enc.Hex);
    console.log(`storedKey: ${storedKey.toString(cryptoJS.enc.Hex)}`);

    var authMessage = `${nonce},${finalNonce},${finalNonce}`;
    var clientProof = scram.clientProof(password, scarmSalt, iter, authMessage);
    console.log(`clientProof: ${clientProof.toString(cryptoJS.enc.Hex)}`);
}

main();