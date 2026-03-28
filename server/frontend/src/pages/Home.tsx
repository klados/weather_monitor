import StatCard from "../components/StatCard.tsx";
import {useEffect, useState} from "react";
import {useTranslation} from "react-i18next";
import ToggleGroup from "../components/ToggleGroup.tsx";

interface WeatherData {
    temperature: number;
    humidity: number;
    recorded_at: string;
}

export function Home() {
    const {t} = useTranslation();

    const [unit, _] = useState("C");
    const [view, setView] = useState("24h");
    const [metric, setMetric] = useState("both");

    const toF = (c: number) => +(c * 9 / 5 + 32).toFixed(1);
    const displayTemp = (c: number) => unit === "C" ? c.toString() : toF(c).toString();
    const tempUnit = unit === "C" ? "°C" : "°F";

    const locationName:string = import.meta.env.VITE_LOCATION_NAME || "Unknown Location";
    const locationCode: string = import.meta.env.VITE_LOCATION_CODE || "";
    const [weatherData, setWeatherData] = useState<WeatherData | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchWeatherData = async () => {
            try {
                setLoading(true);
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:5173';
                const response = await fetch(`${apiUrl}/api/now?location=${encodeURIComponent(locationCode)}`);

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();
                setWeatherData(data);
                setError(null);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to fetch weather data');
                console.error('Error fetching weather data:', err);
            } finally {
                setLoading(false);
            }
        };

        // Call immediately on mount
        fetchWeatherData();
        
        const interval = setInterval(fetchWeatherData, 5*60000);

        return () => clearInterval(interval);
    }, [locationCode]);

    return (
        <div className="relative overflow-x-hidden">

            <h1 className="text-2xl font-bold text-[#f2f3f5] mb-6">{locationName}</h1>

            <div className="grid grid-cols-2 gap-4 mb-4">
                <StatCard
                    label={t("Temperature")}
                    value={loading ? "..." : error ? "—" : displayTemp(weatherData?.temperature || 0)}
                    unit={tempUnit}
                    icon="🌡️"
                    colorClass="text-[#dbdee1]"
                />
                <StatCard
                    label={t("Humidity")}
                    value={loading ? "..." : error ? "—" : weatherData?.humidity.toFixed(1) || "0"}
                    unit="%"
                    icon="💧"
                    colorClass="text-[#dbdee1]"
                />
            </div>

            {error && (
                <div className="mb-4 p-4 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm">
                    {error}
                </div>
            )}

            <div className="rounded-xl p-6 bg-[#2b2d31] shadow-sm">

                {/* Controls */}
                <div className="flex items-center justify-between flex-wrap gap-3 mb-6">
                    <h2 className="text-sm font-bold text-[#f2f3f5] tracking-wide">{t("Historical Data")}</h2>
                    <div className="flex gap-2 flex-wrap">
                        <ToggleGroup
                            options={[["24h", "24h"], ["7d", "7d"]]}
                            value={view} onChange={setView}
                        />
                        <ToggleGroup
                            options={[["both", t("Both")], ["temp", t("Temp")], ["humidity", t("Hum")]]}
                            value={metric} onChange={setMetric}
                        />
                    </div>
                </div>
            </div>

        </div>
    );
}