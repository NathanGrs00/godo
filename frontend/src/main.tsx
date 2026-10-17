import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import './styles/index.css'
import App from './App.tsx'
import AuthProvider from "./context/AuthProvider.tsx";

ReactDOM.createRoot(document.getElementById("root")!).render(
    <BrowserRouter>
        <AuthProvider>
            <App />
        </AuthProvider>
    </BrowserRouter>
);
