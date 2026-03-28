import {useTranslation} from "react-i18next";

export default function About() {
    const {t} = useTranslation();

    return (
        <div className="max-w-4xl mx-auto">
            <h1 className="text-3xl font-bold text-[#f2f3f5] mb-8">
                {t("About Weather Monitor")}
            </h1>

            {/* Main Description */}
            <div className="bg-[#2b2d31] rounded-xl p-6 mb-6 shadow-sm">
                <h2 className="text-xl font-semibold text-[#f2f3f5] mb-4">
                    {t("What is Weather Monitor?")}
                </h2>
                <p className="text-[#dbdee1] leading-relaxed mb-4">
                    {t("AboutDescription")}
                </p>
            </div>

            {/* How It Works */}
            <div className="bg-[#2b2d31] rounded-xl p-6 mb-6 shadow-sm">
                <h2 className="text-xl font-semibold text-[#f2f3f5] mb-4 flex items-center gap-2">
                    <span>⚙️</span> How It Works
                </h2>
                <div className="space-y-4">
                    <div className="flex gap-4">
                        <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center text-white font-bold">
                            1
                        </div>
                        <div>
                            <h3 className="font-semibold text-[#f2f3f5] mb-1">Sensor Collection</h3>
                            <p className="text-[#949ba4] text-sm">
                                ESP32 microcontrollers with DHT sensors collect temperature and humidity data
                                from your environment.
                            </p>
                        </div>
                    </div>

                    <div className="flex gap-4">
                        <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center text-white font-bold">
                            2
                        </div>
                        <div>
                            <h3 className="font-semibold text-[#f2f3f5] mb-1">Backend Processing</h3>
                            <p className="text-[#949ba4] text-sm">
                                A Go-based server receives sensor data and stores it in Firebase, ensuring
                                reliable data persistence and real-time updates.
                            </p>
                        </div>
                    </div>

                    <div className="flex gap-4">
                        <div className="flex-shrink-0 w-8 h-8 rounded-full bg-[#5865f2] flex items-center justify-center text-white font-bold">
                            3
                        </div>
                        <div>
                            <h3 className="font-semibold text-[#f2f3f5] mb-1">Web Visualization</h3>
                            <p className="text-[#949ba4] text-sm">
                                This React-based dashboard displays current conditions and historical trends
                                with interactive charts and real-time updates.
                            </p>
                        </div>
                    </div>
                </div>
            </div>

            {/* Tech Stack */}
            <div className="bg-[#2b2d31] rounded-xl p-6 mb-6 shadow-sm">
                <h2 className="text-xl font-semibold text-[#f2f3f5] mb-4 flex items-center gap-2">
                    <span>🛠️</span> Technology Stack
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="bg-[#1e1f22] rounded-lg p-4 border border-[#40444b]">
                        <h3 className="font-semibold text-[#f2f3f5] mb-2">Hardware</h3>
                        <ul className="text-[#949ba4] text-sm space-y-1">
                            <li>• ESP32 Microcontroller</li>
                            <li>• DHT Temperature/Humidity Sensor</li>
                            <li>• PlatformIO Framework</li>
                        </ul>
                    </div>

                    <div className="bg-[#1e1f22] rounded-lg p-4 border border-[#40444b]">
                        <h3 className="font-semibold text-[#f2f3f5] mb-2">Backend</h3>
                        <ul className="text-[#949ba4] text-sm space-y-1">
                            <li>• Go (Golang)</li>
                            <li>• Firebase / Firestore</li>
                            <li>• RESTful API</li>
                        </ul>
                    </div>

                    <div className="bg-[#1e1f22] rounded-lg p-4 border border-[#40444b]">
                        <h3 className="font-semibold text-[#f2f3f5] mb-2">Frontend</h3>
                        <ul className="text-[#949ba4] text-sm space-y-1">
                            <li>• React + TypeScript</li>
                            <li>• Tailwind CSS</li>
                            <li>• Recharts</li>
                            <li>• Vite</li>
                        </ul>
                    </div>
                </div>
            </div>

            {/* Features */}
            <div className="bg-[#2b2d31] rounded-xl p-6 shadow-sm">
                <h2 className="text-xl font-semibold text-[#f2f3f5] mb-4 flex items-center gap-2">
                    <span>✨</span> Features
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    <div className="flex items-start gap-3">
                        <span className="text-[#3ba55d] text-lg">✓</span>
                        <span className="text-[#dbdee1] text-sm">Real-time temperature & humidity monitoring</span>
                    </div>
                    <div className="flex items-start gap-3">
                        <span className="text-[#3ba55d] text-lg">✓</span>
                        <span className="text-[#dbdee1] text-sm">Historical data visualization (24h / 7 days)</span>
                    </div>
                    <div className="flex items-start gap-3">
                        <span className="text-[#3ba55d] text-lg">✓</span>
                        <span className="text-[#dbdee1] text-sm">Interactive charts with metric filtering</span>
                    </div>
                    <div className="flex items-start gap-3">
                        <span className="text-[#3ba55d] text-lg">✓</span>
                        <span className="text-[#dbdee1] text-sm">Multi-language support (i18n)</span>
                    </div>
                    <div className="flex items-start gap-3">
                        <span className="text-[#3ba55d] text-lg">✓</span>
                        <span className="text-[#dbdee1] text-sm">Responsive design for all devices</span>
                    </div>
                    <div className="flex items-start gap-3">
                        <span className="text-[#3ba55d] text-lg">✓</span>
                        <span className="text-[#dbdee1] text-sm">Automatic data refresh every 5 minutes</span>
                    </div>
                </div>
            </div>
        </div>
    );
}