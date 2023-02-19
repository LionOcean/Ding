/*
 *  @description 判断IP是否为同一局域网
 *  @params localIp 本机IP
 *  @params remoteIp 发送端IP
 * */
export const isEqualLAN = (localIp: string, remoteIp: string): boolean => {

  let localIpArr = localIp.split('.'),
    remoteIpArr = remoteIp.split('.');
  console.log(typeof localIp,localIpArr, remoteIpArr)
  if (localIpArr.length !== remoteIpArr.length) return false;
  let index = 0;
  while (index < localIpArr.length - 1) {
    if (localIpArr[index] !== remoteIpArr[index]) return false;
    index++;
  }
  return true;
};
