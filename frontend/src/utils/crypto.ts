import Crypto from 'crypto-js';

const AES_KEY = 'DING_APP';
/**
 * @param encryptText 需要加密字符串
 * @returns string 加密结果
 */
export const encrypt = (encryptText: string) => Crypto.AES.encrypt(encryptText, AES_KEY).toString();

/**
 * @param decryptText 要解密字符串
 * @returns string 解密结果
 */
export const decrypt = (decryptText: string) => Crypto.AES.decrypt(decryptText, AES_KEY).toString(Crypto.enc.Utf8);
