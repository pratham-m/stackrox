/* eslint-disable */
import React from 'react';
import { Provider } from 'react-redux';
import createHistory from 'history/createBrowserHistory';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { linkTo } from '@storybook/addon-links';

import 'index.css'; // this file is generated by tailwind (see package.json scripts)

import configureStore from 'store/configureStore';
import PanelSlider from 'Components/PanelSlider';

const history = createHistory();
const store = configureStore(undefined, history);

storiesOf('PanelSlider', module)
  .addDecorator(getStory => <Provider store={store}>{getStory()}</Provider>)
  .add('with 3 child elements (panels)', () => {
    return <PanelSlider header="This is a header" onSave={() => action('Saved')}><div className="p-4">Panel 1</div><div className="p-4">Panel 2</div><div className="p-4">Panel 3</div></PanelSlider>;
  })
  .add('with 1 child element (panel)', () => {
    return <PanelSlider header="This is a header" onSave={() => action('Saved')}><div className="p-4">Panel 1</div></PanelSlider>;
  })
  .add('with cancel button', () => {
    return <PanelSlider header="This is a header" onSave={() => action('Saved')} onClose={() => {}}><div className="p-4">Panel 1</div></PanelSlider>;
  });
/* eslint-enable */
