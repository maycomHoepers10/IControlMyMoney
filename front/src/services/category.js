import { getApi } from "./api";

export const createCategory = async (categoryData) => {
    try {
        const api = await getApi();
        const response = await api.post('/category', categoryData);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const updateCategory = async (categoryId, categoryData) => {
    try {
        const api = await getApi();
        const response = await api.put(`/category/${categoryId}`, categoryData);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const deleteCategory = async (categoryId) => {
    try {
        const api = await getApi();
        const response = await api.delete(`/category/${categoryId}`);
        return response.data;
    } catch (error) {
        throw error;
    }
};

export const listCategories = async (setList) => {
    try {
        const api = await getApi();
        const response = await api.get('/categories');

        const data = response.data; // Extrair os dados da resposta

        setList(data); // Atualizar o estado com os dados extraídos
    } catch (error) {
        throw error;
    }
};

export const getCategory = async (categoryId) => {
    try {
        const api = await getApi();
        const response = await api.get(`/category/${categoryId}`);
        const data = response.data;

        return data;
    } catch (error) {
        throw error;
    }
};