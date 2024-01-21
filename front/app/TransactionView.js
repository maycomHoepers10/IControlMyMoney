
import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import { View, Dimensions, StyleSheet, Alert, ActivityIndicator  } from 'react-native';
import { TextInput, Button, Text, HelperText } from 'react-native-paper';
import * as Yup from 'yup';

import { transactionValidationSchema } from "../src/form/validations";
import { createTransaction, getTransaction, updateTransaction } from "../src/services/transaction";

const { width } = Dimensions.get('window');
const rem = width / 16; // Calcula 1rem como 1/16 da largura da tela

export default function TransactionView({ route, navigation }) {
	const [name, setName] = useState(null);
	const [errors, setErrors] = useState({});
	const [loading, setLoading] = useState(true); 
	const { transactionId } = route.params;

	useEffect(() => {

		const fetch = async() => {
			const data = await getTransaction(transactionId);
			setName(data.transactionName);
			setLoading(false);
		};

		if (transactionId) {
			fetch();
		} else {
			setLoading(false);
		}
	}, [transactionId]);

	const handleSubmit = useCallback(async () => {

		try {
			await transactionValidationSchema.validate({ name }, { abortEarly: false });

			if (transactionId) {
				await updateTransaction(transactionId, {
					transactionName: name
				});
			} else {
				await createTransaction({
					transactionName: name
				});
			}

			navigation.navigate('TransactionList');
		} catch (error) {
			if (error instanceof Yup.ValidationError) {
				let validationErrors = {};

				error.inner.forEach((e) => {
					validationErrors[e.path] = e.message;
				});

				setErrors(validationErrors);
			}

			if (axios.isAxiosError(error)) {
				const data = error.response.data;

				Alert.alert('Aviso!', data);
			}
		}
	}, [name, transactionId]);

	if (loading) {
        return (
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
                <ActivityIndicator size="large" color="#0000ff" />
            </View>
        );
    }

	return (
		<View style={{ flex: 1, alignItems: 'center', padding: rem, justifyContent: 'space-between' }}>
			<View style={{ width: '100%' }}>
				<TextInput
					label="Nome"
					value={name}
					onChangeText={text => setName(text)}
				/>
				<HelperText type="error" visible={Boolean(errors.name)}>
					{errors.name}
				</HelperText>
			</View>
			<Button
				icon="content-save"
				mode="contained"
				onPress={handleSubmit}
				style={{ width: '100%' }}
			>
				Salvar
			</Button>
		</View>
	);
}

const styles = StyleSheet.create({
	errorText: {
		color: 'red',
		marginBottom: 8,
	}
});