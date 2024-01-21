import * as React from 'react';
import { View, Text, Button } from 'react-native';
import { NavigationContainer } from '@react-navigation/native';
import {
	createDrawerNavigator,
	DrawerContentScrollView,
	DrawerItemList,
	DrawerItem,
} from '@react-navigation/drawer';

import { createStackNavigator } from '@react-navigation/stack';

import { getJWT } from '../src/services/authentication';

import Summary from './Summary';
import Account from './Account';
import Budget from './Budget';
import CategoryView from './CategoryView';
import CategoryList from './CategoryList';
import TransactionView from './TransactionView';
import TransactionList from './TransactionList';
import Import from './Import';
import Login from './Login';

import { BaseProvider, useBaseContext } from '../src/context/BaseContext';

import {
	Provider as PaperProvider,
	MD3DarkTheme,
	MD3LightTheme,
	MD2DarkTheme,
	MD2LightTheme,
	MD2Theme,
	MD3Theme,
	useTheme,
	adaptNavigationTheme,
} from 'react-native-paper';

function CustomDrawerContent(props) {
	return (
		<DrawerContentScrollView {...props}>
			<DrawerItemList {...props} />
			<DrawerItem
				label="Close drawer"
				onPress={() => props.navigation.closeDrawer()}
			/>
			<DrawerItem
				label="Toggle drawer"
				onPress={() => props.navigation.toggleDrawer()}
			/>
		</DrawerContentScrollView>
	);
}

const Drawer = createDrawerNavigator();
const Stack = createStackNavigator();

function CategoryStack() {
	return (
		<Stack.Navigator screenOptions={{ headerShown: false }}>
			<Stack.Screen name="CategoryList" component={CategoryList} />
			<Stack.Screen name="CategoryView" component={CategoryView} />
		</Stack.Navigator>
	);
}

function TransactionStack() {
	return (
		<Stack.Navigator screenOptions={{ headerShown: false }}>
			<Stack.Screen name="TransactionList" component={TransactionList} />
			<Stack.Screen name="TransactionView" component={TransactionView} />
		</Stack.Navigator>
	);
}

function AuthStack() {
	return (
		<Stack.Navigator screenOptions={{ headerShown: false }}>
			<Stack.Screen name="Login" component={Login} />
		</Stack.Navigator>
	);
}


function MyDrawer() {
	const { jwt } = useBaseContext();

	React.useEffect(() => {
		// Lógica a ser executada quando o JWT é alterado
		console.log('Token JWT alterado:', jwt);
	}, [jwt]);

	if (jwt === undefined) {
		return (
			<View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
				<Text>Verificando o token JWT...</Text>
			</View>
		);
	}

	if (jwt === null) {
		return <AuthStack />;
	}

	console.log(jwt);
	return (
		<Drawer.Navigator
			initialRouteName="Categorias"
		// drawerContent={(props) => <CustomDrawerContent {...props} />}
		>
			<Drawer.Screen name="Resumo" component={Summary} />
			<Drawer.Screen name="Lançamentos" component={TransactionStack} />
			<Drawer.Screen name="Contas" component={Account} initialParams={{ itemId: 42 }} />
			<Drawer.Screen name="Categorias" component={CategoryStack} />
			<Drawer.Screen name="Orçamentos" component={Budget} />
			<Drawer.Screen name="Importar dados" component={Import} />
		</Drawer.Navigator>
	);
}

export default function App() {
	return (
		<PaperProvider>
			<BaseProvider>
				<NavigationContainer independent={true}>
					<MyDrawer />
				</NavigationContainer>
			</BaseProvider>
		</PaperProvider>
	);
}
