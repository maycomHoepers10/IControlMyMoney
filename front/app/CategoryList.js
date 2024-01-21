import * as React from 'react';
import { View, StyleSheet, ScrollView, Alert, TouchableOpacity } from 'react-native';
import { FAB, Portal, Card, Avatar } from 'react-native-paper';
import { useIsFocused } from '@react-navigation/native';
import { listCategories } from "../src/services/category";
import { useFocusEffect } from '@react-navigation/native';
import { deleteCategory } from "../src/services/category";

export default function CategoryList({ navigation }) {

	const [visible, setVisible] = React.useState(false);
	const [open, setOpen] = React.useState(false);
	const isFocused = useIsFocused();
	const [list, setList] = React.useState([]);
	const [selected, setSelected] = React.useState(null);
	const [clicks, setClicks] = React.useState(0);

	React.useEffect(() => {
		setVisible(isFocused);
	}, [isFocused]);

	useFocusEffect(
		React.useCallback(() => {
		  // Carregar os dados ou fazer qualquer operação quando o componente está focado
		  // Por exemplo, você pode chamar a função para carregar dados aqui
		  listCategories(setList);

		  return () => {
			// Função de limpeza, se necessário
		  };
		}, [selected])
	);

	const handleDelete = async() => {
		await deleteCategory(selected);

		setSelected(null);
		navigation.navigate('CategoryList');
	};

	const handleDoublePress = () => {
		setClicks(clicks + 1);
	  
		if (clicks === 1) {
		  // Duplo clique detectado
		  setSelected(null);
		  setClicks(0); // Reseta o contador de cliques
		} else {
		  setTimeout(() => {
			// Redefine o contador de cliques após um intervalo
			setClicks(0);
		  }, 300); // 300 milissegundos é um bom valor para distinguir um clique único de um duplo clique
		}
	  };

	return (
		<View style={styles.container} onPress={handleDoublePress}>
			<TouchableOpacity
				activeOpacity={1}
				style={styles.fullScreen}
				onPress={handleDoublePress}
      		>
				<ScrollView>
					{list.length > 0 && list.map((item) => {
						return (
							<Card
								key={item.id}
								style={styles.card}
								onPress={() => {
									navigation.navigate('CategoryView', {
										categoryId: item.id
									});
								}}
								onLongPress={() => setSelected(selected ? null : item.id)}
								// mode={selectedMode}
							>
								<Card.Title
									title={item.categoryName}
									left={(props) => <Avatar.Icon {...props} icon="plus" />}
								/>
							</Card>
						);
					})}
				</ScrollView>
				<View>
					<Portal>
						<FAB.Group
							open={false}
							icon={selected ? 'trash-can' : 'plus'}
							actions={[
								{ icon: 'plus', label: 'Incluir', onPress: () => { } }
							]}
							onStateChange={({ open }) => { setOpen(open) }}
							onPress={() => {
								if (open) {

									if (selected) {
										Alert.alert(
											'Confirmação',
											'Tem certeza que deseja excluir esse registro?',
											[
											{
												text: 'Cancelar',
												onPress: () => {
													setSelected(null);
												},
												style: 'cancel',
											},
											{
												text: 'Confirmar',
												onPress: () => {
													handleDelete();
												},
											},
											],
											{ cancelable: false }
										);
									} else {
										navigation.navigate('CategoryView', { categoryId: null });
									}
								}
							}}
							visible={visible}
						/>
					</Portal>
				</View>
			</TouchableOpacity>
		</View>
	);
}



const styles = StyleSheet.create({
	container: {
	  flex: 1,
	  backgroundColor: '#F8F8FF'
	},
	content: {
	  padding: 4,
	},
	card: {
	  margin: 10,
	  backgroundColor: '#13395D'
	},
	chip: {
	  margin: 4,
	},
	preference: {
	  alignItems: 'center',
	  flexDirection: 'row',
	  paddingVertical: 12,
	  paddingHorizontal: 8,
	},
	button: {
	  borderRadius: 12,
	},
	fullScreen: {
		...StyleSheet.absoluteFillObject,
	},
  });
