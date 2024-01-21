import * as React from 'react';
import { View, StyleSheet, ScrollView } from 'react-native';
import { FAB, Portal, useTheme, Text, Card, Provider, Avatar } from 'react-native-paper';
import { useIsFocused } from '@react-navigation/native';


export default function Account() {

	const [visible, setVisible] = React.useState(false);
	const [open, setOpen] = React.useState(false);
	const isFocused = useIsFocused();

	React.useEffect(() => {
		setVisible(isFocused);
	}, [isFocused]);


	return (
		<View style={styles.container}>
			<ScrollView>
			<Card
				style={styles.card}
				// onPress={() => {
				// 	preferences.toggleTheme();
				// }}
				// mode={selectedMode}
			>
				<Card.Title
					title="Pressable Theme Change"
					left={(props) => <Avatar.Icon {...props} icon="plus" />}
				/>
			</Card>
			</ScrollView>
			<View>
				<Portal>
					<FAB.Group
						open={false}
						icon={'plus'}
						actions={[
							{ icon: 'plus', label: 'Incluir', onPress: () => { } },
							{ icon: 'star', label: 'Star', onPress: () => { } }
						]}
						onStateChange={({ open }) => { setOpen(open) }}
						onPress={() => {
							if (open) {
								// Faça algo se o FAB estiver aberto
							}
						}}
						visible={visible}
					/>
				</Portal>
			</View>
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
  });