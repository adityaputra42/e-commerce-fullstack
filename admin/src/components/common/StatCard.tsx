interface StatCardProps {
  title: string;
  value: string | number;
  icon: React.ReactNode;
  color: 'blue' | 'green' | 'yellow' | 'red';
}

const colorMap = {
  blue: 'bg-blue-100 text-blue-600',
  green: 'bg-green-100 text-green-600',
  yellow: 'bg-yellow-100 text-yellow-600',
  red: 'bg-red-100 text-red-600',
};

export const StatCard = ({
  title,
  value,
  icon,
  color,
}: StatCardProps) => (
  <div className="bg-white p-6 rounded-lg shadow-md flex justify-between items-center">
    <div>
      <p className="text-sm text-gray-500">{title}</p>
      <p className="text-3xl font-bold text-gray-900">{value}</p>
    </div>
    <div
      className={`p-3 rounded-full ${colorMap[color]}`}
    >
      {icon}
    </div>
  </div>
);
