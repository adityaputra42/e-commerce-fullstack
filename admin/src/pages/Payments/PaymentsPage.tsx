import { useState, useEffect } from 'react';
import api from '../../services/api';
import { Search, Filter, ArrowUpRight, CheckCircle2, Clock, XCircle, MoreVertical } from 'lucide-react';

interface Payment {
  id: number;
  transaction_id: string;
  total_payment: number;
  status: string;
  created_at: string;
  transaction?: {
    tx_id: string;
    customer_name?: string;
  };
}

const PaymentsPage = () => {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetchPayments();
  }, []);

  const fetchPayments = async () => {
    setIsLoading(true);
    try {
      const response = await api.get('/payments');
      // Backend returns data in a "data" property
      const data = response.data?.data || [];
      setPayments(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Failed to fetch:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const getStatusStyle = (status: string) => {
    switch (status.toLowerCase()) {
      case 'success':
      case 'confirmed':
      case 'completed':
        return 'bg-emerald-50 text-emerald-600 ring-emerald-100';
      case 'pending':
        return 'bg-amber-50 text-amber-600 ring-amber-100';
      case 'failed':
      case 'rejected':
      case 'cancelled':
        return 'bg-rose-50 text-rose-600 ring-rose-100';
      default:
        return 'bg-slate-50 text-slate-600 ring-slate-100';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'success':
      case 'confirmed':
      case 'completed':
        return <CheckCircle2 className="w-3.5 h-3.5" />;
      case 'pending':
        return <Clock className="w-3.5 h-3.5" />;
      default:
        return <XCircle className="w-3.5 h-3.5" />;
    }
  };

  const filteredPayments = payments.filter(p => 
     p.transaction_id.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-extrabold text-slate-900 tracking-tight transition-all">Payments History</h1>
          <p className="text-slate-500 mt-1 font-medium">Monitor and manage all incoming customer transactions.</p>
        </div>
        
        <div className="flex items-center gap-3">
             <div className="relative group">
                <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 group-focus-within:text-indigo-600 transition-colors" />
                <input 
                    type="text" 
                    placeholder="Search by ID..." 
                    className="bg-white border-2 border-slate-50 rounded-2xl py-2.5 pl-10 pr-4 text-sm focus:border-indigo-500/20 focus:ring-4 focus:ring-indigo-500/10 outline-none transition-all w-64 font-medium"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                />
             </div>
             <button className="p-2.5 bg-white border-2 border-slate-50 rounded-2xl text-slate-600 hover:bg-slate-50 transition-colors">
                <Filter className="w-5 h-5" />
             </button>
        </div>
      </div>

      <div className="premium-card overflow-hidden">
        {isLoading ? (
            <div className="p-20 flex flex-col items-center justify-center gap-4">
                <div className="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin"></div>
                <p className="text-slate-400 font-bold text-xs uppercase tracking-widest">Verifying ledger...</p>
            </div>
        ) : (
            <div className="overflow-x-auto">
                <table className="w-full">
                    <thead>
                        <tr className="bg-slate-50/50 border-b border-slate-100">
                            <th className="px-6 py-4 text-left text-[10px] font-black text-slate-400 uppercase tracking-widest">Transaction Info</th>
                            <th className="px-6 py-4 text-left text-[10px] font-black text-slate-400 uppercase tracking-widest">Date & Time</th>
                            <th className="px-6 py-4 text-left text-[10px] font-black text-slate-400 uppercase tracking-widest">Total Amount</th>
                            <th className="px-6 py-4 text-left text-[10px] font-black text-slate-400 uppercase tracking-widest">Status</th>
                            <th className="px-6 py-4 text-right text-[10px] font-black text-slate-400 uppercase tracking-widest">Actions</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-50">
                        {filteredPayments.map((payment) => (
                            <tr key={payment.id} className="group hover:bg-slate-50/50 transition-colors">
                                <td className="px-6 py-5 whitespace-nowrap">
                                    <div className="flex items-center gap-3">
                                        <div className="w-10 h-10 bg-white rounded-xl shadow-sm border border-slate-100 flex items-center justify-center group-hover:scale-110 transition-transform">
                                            <ArrowUpRight className="w-5 h-5 text-indigo-500" />
                                        </div>
                                        <div>
                                            <p className="text-sm font-black text-slate-900 font-mono tracking-tighter">#{payment.transaction_id}</p>
                                            <p className="text-[10px] text-slate-400 font-bold uppercase">External Reference</p>
                                        </div>
                                    </div>
                                </td>
                                <td className="px-6 py-5 whitespace-nowrap">
                                    <p className="text-sm font-bold text-slate-700">
                                        {new Date(payment.created_at).toLocaleDateString('en-US', { day: 'numeric', month: 'short', year: 'numeric' })}
                                    </p>
                                    <p className="text-[10px] text-slate-400 font-medium">
                                        {new Date(payment.created_at).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' })}
                                    </p>
                                </td>
                                <td className="px-6 py-5 whitespace-nowrap">
                                    <span className="text-sm font-black text-slate-900 italic">
                                        {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(payment.total_payment)}
                                    </span>
                                </td>
                                <td className="px-6 py-5 whitespace-nowrap">
                                    <span className={`inline-flex items-center gap-1.5 px-3 py-1 ring-1 ring-inset rounded-full text-[10px] font-black uppercase tracking-widest ${getStatusStyle(payment.status)}`}>
                                        {getStatusIcon(payment.status)}
                                        {payment.status}
                                    </span>
                                </td>
                                <td className="px-6 py-5 whitespace-nowrap text-right">
                                    <button className="p-2 text-slate-400 hover:text-indigo-600 hover:bg-white rounded-xl transition-all shadow-sm">
                                        <MoreVertical className="w-4 h-4" />
                                    </button>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
                {filteredPayments.length === 0 && (
                    <div className="p-20 text-center">
                        <p className="text-slate-400 font-medium italic">No payment records found.</p>
                    </div>
                )}
            </div>
        )}
      </div>
    </div>
  );
};

export default PaymentsPage;
