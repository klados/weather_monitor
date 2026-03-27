type StatCardProps = {
    label: string;
    value: string;
    unit: string;
    colorClass: string;
    icon: string;
};

function StatCard({ label, value, unit, colorClass, icon } : StatCardProps) {
    return (
        <div className="rounded-xl p-5 bg-[#2b2d31] shadow-sm">
            <div className="flex items-start justify-between mb-3">
                <span className="text-2xl">{icon}</span>
                <span className={`text-xs font-bold tracking-wide ${colorClass}`}>
                    {label}
                </span>
            </div>
            <div className="flex items-end gap-1">
                <span className="text-4xl font-semibold text-[#f2f3f5]">
                  {value}
                </span>
                <span className="text-lg mb-1 text-[#b5bac1]">{unit}</span>
            </div>
        </div>
    );
}

export default StatCard;