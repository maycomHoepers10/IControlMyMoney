import * as SecureStore from 'expo-secure-store';
import api from './api';

// Para salvar o token
const saveToken = async (token) => {
  await SecureStore.setItemAsync('jwtToken', token);
};


// Para recuperar o token
export const getJWT = async () => {
    const jwtToken = await SecureStore.getItemAsync('jwtToken');

    if (jwtToken) {
        const response = await api.post('/validate-token', {
            token: jwtToken
        });
    
        if (response.status === 200) {
            const { isValid } = response.data;

            return isValid ? jwtToken : null;
        }
    }

    return null;
};  

export const authentication = async (data) => {
    try {
        const response = await api.post('/users/login', data);

        if (response.status === 200) {
            saveToken(response.data.token);
        }
    } catch (error) {
        throw error;
    }
};
