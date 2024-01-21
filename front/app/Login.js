import React, { useState } from 'react';
import { View, StyleSheet } from 'react-native';
import { TextInput, Button } from 'react-native-paper';
import { authentication } from '../src/services/authentication';
import * as SecureStore from 'expo-secure-store';
import { useBaseContext } from '../src/context/BaseContext';

export default function Login() {
	const { setJwt } = useBaseContext();
	const [username, setUsername] = useState('');
  	const [password, setPassword] = useState('');

	const handleLogin = async() => {
		await authentication({ Email: username, Password: password });
		const token = await SecureStore.getItemAsync('jwtToken');

		setJwt(token);
	};

	return (
		<View style={styles.container}>
			<TextInput
				label="Username"
				value={username}
				onChangeText={(text) => setUsername(text)}
				style={styles.input}
			/>
			<TextInput
				label="Password"
				value={password}
				onChangeText={(text) => setPassword(text)}
				secureTextEntry
				style={styles.input}
			/>
			<Button mode="contained" onPress={handleLogin} style={styles.button}>
				Entrar
			</Button>
		</View>
	);
}

const styles = StyleSheet.create({
	container: {
		flex: 1,
		justifyContent: 'center',
		padding: 16,
	},
	input: {
	  	marginBottom: 16,
	},
	button: {
	  	marginTop: 16,
	},
});