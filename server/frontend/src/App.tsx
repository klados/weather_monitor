import './App.css'
import { Link, Outlet } from "react-router";
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import enTranslations from "./locales/en.json";
import elTranslations from "./locales/el.json";
import { useTranslation } from "react-i18next";


i18n
    .use(initReactI18next) // passes i18n down to react-i18next
    .init({
        resources: {
            en: {
                translation: enTranslations
            },
            el: {
                translation: elTranslations
            }
        },
        lng: "en",
        fallbackLng: "en",

        interpolation: {
            escapeValue: false
        }
    });



function App() {

    const { t, i18n } = useTranslation();

    const changeLanguage = (lng: string) => {
        i18n.changeLanguage(lng);
    };

    return (
      <div className="min-h-screen bg-[#313338] text-[#dbdee1] font-sans">
          <nav className="p-4 mb-4 flex justify-between items-center bg-[#1e1f22] shadow-sm">
              <ul className="flex gap-4 list-none m-0 p-0">
                  <li>
                      <Link to="/" className="text-[#dbdee1] hover:text-white font-medium transition-colors">
                          {t('HomePage')}
                      </Link>
                  </li>
                  <li>
                      <Link to="/about" className="text-[#dbdee1] hover:text-white font-medium transition-colors">
                          {t('AboutPage')}
                      </Link>
                  </li>
              </ul>
              <div className="flex gap-2">
                  <button
                      onClick={() => changeLanguage('en')}
                      className={`px-3 py-1 text-sm font-medium rounded transition-colors ${i18n.language === 'en' ? 'bg-[#5865F2] text-white' : 'bg-[#2b2d31] text-[#dbdee1] hover:bg-[#3f4147]'}`}
                  >
                      EN
                  </button>
                  <button
                      onClick={() => changeLanguage('el')}
                      className={`px-3 py-1 text-sm font-medium rounded transition-colors ${i18n.language === 'el' ? 'bg-[#5865F2] text-white' : 'bg-[#2b2d31] text-[#dbdee1] hover:bg-[#3f4147]'}`}
                  >
                      ΕΛ
                  </button>
              </div>
          </nav>

          <main className="p-4">
              <Outlet />
          </main>
      </div>
  )
}

export default App
