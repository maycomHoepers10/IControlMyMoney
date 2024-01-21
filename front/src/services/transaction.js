import { getApi } from "./api";

export const createTransaction = async (transactionData) => {
    try {
        const api = await getApi();
        const response = await api.post('/transaction', transactionData);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const updateTransaction = async (transactionId, transactionData) => {
    try {
        const api = await getApi();
        const response = await api.put(`/transaction/${transactionId}`, transactionData);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const deleteTransaction = async (transactionId) => {
    try {
        const api = await getApi();
        const response = await api.delete(`/transaction/${transactionId}`);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const listTransactions = async (setList) => {
    try {
        const api = await getApi();
        const response = await api.get('/transactions');

        const data = response.data; // Extrair os dados da resposta

        setList(data); // Atualizar o estado com os dados extraídos
    } catch (error) {
        throw error;
    }
};

export const getTransaction = async (transactionId) => {
    try {
        const api = await getApi();
        const response = await api.get(`/transaction/${transactionId}`);
        const data = response.data;

        return data;
    } catch (error) {
        throw error;
    }
};