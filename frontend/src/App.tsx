import { HashRouter } from 'react-router-dom';
import Router from './routers';

import { Layout } from 'antd';

import Header from './components/Header';
import Navigate from './components/Navigate';

import './App.less';

const { Content, Footer } = Layout;

function App() {
  return (
    <HashRouter>
      <Layout className='app__layout'>
        <Header />
        <Content>
          <Router />
        </Content>
        <Navigate />
      </Layout>
    </HashRouter>
  );
}

export default App;
