import StatCard from "../components/StatCard.tsx";
import {useState} from "react";
import {useTranslation} from "react-i18next";
import ToggleGroup from "../components/ToggleGroup.tsx";

export function Home() {
    const {t} = useTranslation();

    const [unit, _] = useState("C");
    const [view, setView] = useState("24h");
    const [metric, setMetric] = useState("both");

    const toF = (c: number) => +(c * 9 / 5 + 32).toFixed(1);
    const displayTemp = (c: number) => unit === "C" ? c.toString() : toF(c).toString();
    const tempUnit = unit === "C" ? "°C" : "°F";

    return (
        <div className="relative overflow-x-hidden">

            <div className="grid grid-cols-2 gap-4 mb-4">
                <StatCard
                    label={t("Temperature")} value={displayTemp(12)}
                    unit={tempUnit} icon="🌡️" colorClass="text-[#dbdee1]"
                />
                <StatCard
                    label={t("Humidity")} value={"0"}
                    unit="%" icon="💧" colorClass="text-[#dbdee1]"
                />
            </div>

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