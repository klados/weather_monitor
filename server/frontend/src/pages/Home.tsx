import StatCard from "../components/StatCard.tsx";
import {useEffect, useState} from "react";
import {useTranslation} from "react-i18next";
import ToggleGroup from "../components/ToggleGroup.tsx";
import {LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer} from 'recharts';

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
    const [historicalData, setHistoricalData] = useState<WeatherData[]>([]);
    const [loading, setLoading] = useState(true);
    const [historicalLoading, setHistoricalLoading] = useState(true);
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


    useEffect(() => {
        const fetchHistoricalData = async () => {
            try {
                setHistoricalLoading(true);
                const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:5173';
                const timespanInDays = view === "24h" ? 1 : 7;
                const response = await fetch(`${apiUrl}/api/historicalData?location=${encodeURIComponent(locationCode)}&timespanInDays=${timespanInDays}`);

                if (!response.ok) {
                    throw new Error(`HTTP error! status: ${response.status}`);
                }

                const data = await response.json();
                setHistoricalData(data);
            } catch (err) {
                console.error('Error fetching historical data:', err);
            } finally {
                setHistoricalLoading(false);
            }
        };

        fetchHistoricalData();
    }, [locationCode, view]);

    const chartData = historicalData.map(item => ({
        time: new Date(item.recorded_at).toLocaleString('en-US', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }),
        temperature: unit === "C" ? item.temperature : toF(item.temperature),
        humidity: item.humidity
    }));


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

                {historicalLoading ? (
                    <div className="h-80 flex items-center justify-center text-[#949ba4]">
                        {t("Loading...")}
                    </div>
                ) : (
                    <ResponsiveContainer width="100%" height={320}>
                        <LineChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3" stroke="#40444b" />
                            <XAxis
                                dataKey="time"
                                stroke="#949ba4"
                                tick={{fill: '#949ba4', fontSize: 12}}
                            />
                            <YAxis stroke="#949ba4" tick={{fill: '#949ba4', fontSize: 12}} />
                            <Tooltip
                                contentStyle={{
                                    backgroundColor: '#1e1f22',
                                    border: '1px solid #40444b',
                                    borderRadius: '8px',
                                    color: '#f2f3f5'
                                }}
                            />
                            <Legend wrapperStyle={{color: '#f2f3f5'}} />
                            {(metric === "both" || metric === "temp") && (
                                <Line
                                    type="monotone"
                                    dataKey="temperature"
                                    stroke="#ed4245"
                                    name={`${t("Temperature")} (${tempUnit})`}
                                    strokeWidth={2}
                                    dot={false}
                                />
                            )}
                            {(metric === "both" || metric === "humidity") && (
                                <Line
                                    type="monotone"
                                    dataKey="humidity"
                                    stroke="#5865f2"
                                    name={`${t("Humidity")} (%)`}
                                    strokeWidth={2}
                                    dot={false}
                                />
                            )}
                        </LineChart>
                    </ResponsiveContainer>
                )}
            </div>

        </div>
    );
}