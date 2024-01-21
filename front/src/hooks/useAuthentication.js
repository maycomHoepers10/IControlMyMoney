import { useState, useEffect } from 'react';
import { getJWT } from '../services/authentication';

const useAuthentication = () => {
	const [jwt, setJwt] = useState(undefined);

	useEffect(() => {
		const fetchJWT = async () => {
			const token = await getJWT();
			setJwt(token);
		};

		fetchJWT();
	}, []);

	return { jwt, setJwt };
};

export default useAuthentication;