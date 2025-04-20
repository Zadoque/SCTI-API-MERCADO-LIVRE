import Head from 'next/head';
import PaymentForm from '../components/PaymentForm';

export default function Home() {
  return (
    <div className="container mx-auto py-8">
      <Head>
        <title>Teste de Pagamento - Mercado Pago</title>
        <meta name="description" content="Teste de integração com Mercado Pago" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main className="max-w-4xl mx-auto">
        <h1 className="text-3xl font-bold text-center mb-8">Teste de Pagamento</h1>
        
        <PaymentForm />
      </main>
    </div>
  );
}
