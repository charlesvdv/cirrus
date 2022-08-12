/* @refresh reload */
import { render } from 'solid-js/web'
import { Router } from "@solidjs/router";

import App from './App'
import './assets/index.css'
import '@fontsource/inter/variable-full.css'

render(
  () => (
    <Router>
      <App />
    </Router>
  ),
  document.getElementById('root') as HTMLElement
)
