import { createApp } from "vue";
import "./style.css";
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";
import { library } from "@fortawesome/fontawesome-svg-core";
import { CustomButton } from "./components/shared";
import {
  faArrowRight,
  faGear,
  faDatabase,
  faFileCode,
  faLock,
  faArrowLeft,
  faEye,
  faEyeSlash,
  faCancel,
  faX,
  faSpinner,
  faAdd,
  faExternalLink,
  faTrashCan,
  faPencilAlt,
  faBars,
  faWarning,
  faClipboard,
} from "@fortawesome/free-solid-svg-icons";
import { faCheckCircle } from "@fortawesome/free-regular-svg-icons";
import {
  faGithub,
  faGoogle,
  faTwitter,
  faWhatsapp,
} from "@fortawesome/free-brands-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
import Toast, { PluginOptions } from "vue-toastification";
// Import the CSS or use your own!
import "vue-toastification/dist/index.css";

// pinia - for global state management
const pinia = createPinia();

const toastOptions: PluginOptions = {};

library.add(
  faArrowRight,
  faArrowLeft,
  faGear,
  faDatabase,
  faFileCode,
  faLock,
  faGithub,
  faGoogle,
  faTwitter,
  faWhatsapp,
  faEye,
  faEyeSlash,
  faCancel,
  faX,
  faSpinner,
  faAdd,
  faExternalLink,
  faTrashCan,
  faPencilAlt,
  faClipboard,
  faBars,
  faWarning,
  faCheckCircle
);

createApp(App)
  .use(router) // router
  .use(pinia)
  .use(Toast, toastOptions) // toast notifications
  .component("font-awesome-icon", FontAwesomeIcon) // font awesome icons access
  .component("custom-button", CustomButton)
  .mount("#app");
