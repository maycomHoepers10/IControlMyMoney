import * as SecureStore from 'expo-secure-store';
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://192.168.100.5:54165/',
  headers: {
	'Content-Type': 'application/json',
  },
});


export const getApi = async () => {

	const token = await SecureStore.getItemAsync('jwtToken');
	
	if (token) {
		// Adicione o token Bearer ao cabeçalho da requisição
		api.defaults.headers.common['Authorization'] = `Bearer ${token}`;
	}

	return api;
};

export default api;