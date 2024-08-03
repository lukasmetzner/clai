import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import PrimeVue from 'primevue/config';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Aura from '@primevue/themes/aura';
import Textarea from 'primevue/textarea';
import '/node_modules/primeflex/primeflex.css'


const app = createApp(App);
app.use(PrimeVue, {
    theme: {
        preset: Aura
    }
});
app.component('InputText', InputText);
app.component('Button', Button);
app.component('TextArea',  Textarea);
app.mount('#app');
