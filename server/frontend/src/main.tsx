import { createRoot } from 'react-dom/client'
import './index.css'
import { BrowserRouter, Routes, Route } from "react-router";
import Home from "./pages/Home.tsx";
import About from "./pages/About.tsx";
import App from './App.tsx';

createRoot(document.getElementById('root')!).render(
    <BrowserRouter>
        <Routes>
            <Route element={<App />}>
                <Route index element={<Home />} />
                <Route path="about" element={<About />} />
            </Route>
        </Routes>
    </BrowserRouter>,
)
