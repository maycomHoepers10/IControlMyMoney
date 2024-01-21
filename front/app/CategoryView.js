
import React, { useEffect, useState, useCallback } from 'react';
import axios from 'axios';
import { View, Dimensions, StyleSheet, Alert, ActivityIndicator  } from 'react-native';
import { TextInput, Button, Text, HelperText } from 'react-native-paper';
import * as Yup from 'yup';

import { categoryValidationSchema } from "../src/form/validations";
import { createCategory, getCategory, updateCategory } from "../src/services/category";

const { width } = Dimensions.get('window');
const rem = width / 16; // Calcula 1rem como 1/16 da largura da tela

export default function CategoryView({ route, navigation }) {
	const [name, setName] = useState(null);
	const [errors, setErrors] = useState({});
	const [loading, setLoading] = useState(true); 
	const { categoryId } = route.params;

	useEffect(() => {

		const fetch = async() => {
			const data = await getCategory(categoryId);
			setName(data.categoryName);
			setLoading(false);
		};

		if (categoryId) {
			fetch();
		} else {
			setLoading(false);
		}
	}, [categoryId]);

	const handleSubmit = useCallback(async () => {

		try {
			await categoryValidationSchema.validate({ name }, { abortEarly: false });

			if (categoryId) {
				await updateCategory(categoryId, {
					categoryName: name
				});
			} else {
				await createCategory({
					categoryName: name
				});
			}

			navigation.navigate('CategoryList');
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
	}, [name, categoryId]);

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