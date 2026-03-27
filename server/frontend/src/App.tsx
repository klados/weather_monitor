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
      <>
          <nav className="p-4 border-b border-gray-300 mb-4 flex justify-between items-center bg-white shadow-sm">
              <ul className="flex gap-4 list-none m-0 p-0">
                  <li>
                      <Link to="/" className="text-blue-600 hover:text-blue-800 font-medium transition-colors">
                          {t('HomePage')}
                      </Link>
                  </li>
                  <li>
                      <Link to="/about" className="text-blue-600 hover:text-blue-800 font-medium transition-colors">
                          {t('AboutPage')}
                      </Link>
                  </li>
              </ul>
              <div className="flex gap-2">
                  <button
                      onClick={() => changeLanguage('en')}
                      className={`px-3 py-1 text-sm font-semibold rounded transition-colors ${i18n.language === 'en' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}`}
                  >
                      EN
                  </button>
                  <button
                      onClick={() => changeLanguage('el')}
                      className={`px-3 py-1 text-sm font-semibold rounded transition-colors ${i18n.language === 'el' ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}`}
                  >
                      ΕΛ
                  </button>
              </div>
          </nav>


          <main>
              <Outlet />
          </main>
      </>
  )
}

export default App
