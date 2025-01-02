const CryptoJS = require('crypto-js');

function computeHMAC(keyHex, message) {
    const key = CryptoJS.enc.Hex.parse(keyHex);
    const hmac = CryptoJS.HmacSHA256(message, key);
    return hmac.toString(CryptoJS.enc.Hex);
}

// Example Usage
const keyHex = '4ec91c5b362986cec5ef6db10eb4bb5b1c54eef1eff89006cf00a8993226fabd';
const message = 'Client Key';
const hmacResult = computeHMAC(keyHex, message);
console.log('HMAC Result:', hmacResult);