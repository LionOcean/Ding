/*
 *  @description 判断IP是否为同一局域网
 *  @params localIp 本机IP
 *  @params remoteIp 发送端IP
 *  @params mask 子网掩码
 * */
export const isEqualLAN = (ip1: string, ip2: string, mask: string = '255.255.252.0'): boolean => {
  const ip1Arr = ip1.split('.');
  const ip2Arr = ip2.split('.');
  const maskArr = mask.split('.');
  console.log(ip1, ip2);
  for (let i = 0; i < 4; i++) {
    if ((parseInt(ip1Arr[i]) & parseInt(maskArr[i])) !== (parseInt(ip2Arr[i]) & parseInt(maskArr[i]))) {
      return false;
    }
  }
  return true;
};

type byteUnit = 'b' | 'Kb' | 'Mb' | 'Gb';

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

  while (byteSize >= base) {
    byteSize = byteSize / base;
    i++;
  }

  return [Number(byteSize.toFixed(2)), units[i]];
};
