import { getApi } from "./api";

export const createAccount = async (accountData) => {
    try {
        const api = await getApi();
        const response = await api.post('/financialAccount', accountData);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const updateAccount = async (accountId, accountData) => {
    try {
        const api = await getApi();
        const response = await api.put(`/financialAccount/${accountId}`, accountData);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const deleteAccount = async (accountId) => {
    try {
        const api = await getApi();
        const response = await api.delete(`/financialAccount/${accountId}`);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const listAccounts = async (setList) => {
    try {
        const api = await getApi();
        const response = await api.get('/financialAccounts');

        const data = response.data; // Extrair os dados da resposta

        setList(data); // Atualizar o estado com os dados extraídos
    } catch (error) {
        throw error;
    }
};

export const getAccount = async (accountId) => {
    try {
        const api = await getApi();
        const response = await api.get(`/financialAccount/${accountId}`);
        const data = response.data;

        return data;
    } catch (error) {
        throw error;
    }
};