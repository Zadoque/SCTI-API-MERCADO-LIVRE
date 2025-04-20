'use client'
import React, { useState } from 'react';
import { initMercadoPago, CardPayment } from '@mercadopago/sdk-react';

await initMercadoPago(''); //CHAVE PÚBLICA AQUI 

export default function PaymentForm() {
  const [loading, setLoading] = useState(false);
  const [paymentResult, setPaymentResult] = useState(null);

  const initialization = {
    amount: 1, 
  };

  const customization = {
    visual: {
      hidePaymentButton: false, 
      buttonText: 'Pagar agora', 
      buttonBackground: '#0046c0', 
    },
    paymentMethods: {
      maxInstallments: 1,
    },
  };

  const onSubmit = async (formData) => {
    setLoading(true);
    console.log("Dados do formulário recebidos:", formData);
    
    try {
      const response = await fetch('http://localhost:8080/api/payment', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token: formData.token,
          amount: initialization.amount,
          description: 'Compra de produto',
          email: formData.payer.email,
          payment_method: formData.payment_method_id,
        }),
      });
      
      const result = await response.json();
      setPaymentResult(result);
      console.log("Resposta do servidor:", result);
    } catch (error) {
      console.error('Erro ao processar pagamento:', error);
      setPaymentResult({
        success: false,
        message: 'Erro ao processar pagamento',
      });
    } finally {
      setLoading(false);
    }
  };

  const onError = (error) => {
    console.error('Erro no formulário de pagamento:', error);
  };

  const onReady = () => {
    console.log('Formulário pronto');
  };

  return (
    <div className="max-w-md mx-auto p-4 bg-white rounded-lg shadow-md">
      <h2 className="text-xl font-bold mb-4">Pagamento</h2>
      
      <div className="mb-4">
        <p className="text-gray-700 mb-2">Valor: R$ {initialization.amount.toFixed(2)}</p>
        <p className="text-gray-700 mb-2">Produto: Exemplo de Produto</p>
      </div>
      
      <div className="mb-4 border p-4 rounded-lg">
        <CardPayment
          initialization={initialization}
          customization={customization}
          onSubmit={onSubmit}
          onReady={onReady}
          onError={onError}
        />
      </div>
      
      {loading && (
        <div className="mt-4 p-3 bg-blue-50 rounded-lg">
          <p className="text-blue-700">Processando pagamento...</p>
        </div>
      )}
      
      {paymentResult && (
        <div className={`mt-4 p-3 rounded-lg ${paymentResult.success ? 'bg-green-100' : 'bg-red-100'}`}>
          <p className="font-medium">{paymentResult.message}</p>
          {paymentResult.success && <p className="text-sm mt-1">Pagamento processado com sucesso!</p>}
        </div>
      )}
    </div>
  );
}
