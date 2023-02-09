import { WindowMinimise, Quit } from '../../../wailsjs/runtime';
import './index.less';

export default function () {
  return (
    <div
      className='header'
      style={{
        ['--wails-draggable' as string]: 'drag',
      }}
    >
      <div className='header__backBtn' onClick={() => window.history.back()} role='button' tabIndex={-1}>
        返回
      </div>
      <span className='header__title'>我是拖动头部</span>
      <div className='header__operation'>
        <span className='header__minimizeBtn' onClick={() => WindowMinimise()} role='button' tabIndex={-1}>
          最小化
        </span>
        <span className='header__closeBtn' onClick={() => Quit()} role='button' tabIndex={-1}>
          退出
        </span>
      </div>
    </div>
  );
}
