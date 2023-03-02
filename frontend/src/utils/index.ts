/*
 *  @description 判断IP是否为同一局域网
 *  @params localIp 本机IP
 *  @params remoteIp 发送端IP
 * */
export const isEqualLAN = (localIp: string, remoteIp: string): boolean => {
  let localIpArr = localIp.split('.'),
    remoteIpArr = remoteIp.split('.');
  console.log(typeof localIp, localIpArr, remoteIpArr);
  if (localIpArr.length !== remoteIpArr.length) return false;
  let index = 0;
  while (index < localIpArr.length - 1) {
    if (localIpArr[index] !== remoteIpArr[index]) return false;
    index++;
  }
  return true;
};

type byteUnit = 'b'|'Kb'|'Mb'|'Gb';

/**
 * 计算字节大小单位并格式化，结果保留两位小数
 * @param byteSize 字节大小
 * @returns 
 */
export const calcByteUnit = (byteSize: number): [number, byteUnit] => {
  const base = 1024;
  const units: byteUnit[] = ['b', 'Kb', 'Mb', 'Gb'];
  
  let i = 0;
  
  if (byteSize < base) {
    return [byteSize, units[i]];
  }

  while(byteSize >= base) {
    byteSize = byteSize / base
    i++
  }

  
  return [Number(byteSize.toFixed(2)), units[i]]
}