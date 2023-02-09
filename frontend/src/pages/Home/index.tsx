import logo from '../../assets/images/logo.png';

import './index.less';

function IndexPage() {
  return (
    <div className='home'>
      <img className='home__logo' src={logo} alt='logo' />
      <title>Ding</title>
      <div className={'description'}>咻的一下就过去了</div>
    </div>
  );
}

export default IndexPage;
