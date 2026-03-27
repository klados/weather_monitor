type ToggleGroupProps = {
    options: [string, string][],
    value: string,
    onChange: (v: string) => void
};

function ToggleGroup({ options, value, onChange } : ToggleGroupProps) {
    return (
        <div className="flex gap-1 p-1 bg-[#1e1f22] rounded-lg">
            {options.map(([v, l]) => (
                <button
                    key={v}
                    onClick={() => onChange(v)}
                    className={`px-3 py-1.5 rounded-md text-sm font-medium border-none cursor-pointer transition-colors
                ${value === v ? "bg-[#5865F2] text-white shadow-sm" : "bg-transparent text-[#b5bac1] hover:bg-[#313338] hover:text-[#dbdee1]"}`}
                >
                    {l}
                </button>
            ))}
        </div>
    );
}

export default ToggleGroup;