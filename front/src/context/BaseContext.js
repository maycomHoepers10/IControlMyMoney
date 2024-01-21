import { createContext, useContext } from 'react';
import useAuthentication from '../hooks/useAuthentication'

const BaseContext = createContext();


export const BaseProvider = ({ children }) => {
  const { jwt, setJwt } = useAuthentication();

  return (
    <BaseContext.Provider value={{ jwt, setJwt }}>
      {children}
    </BaseContext.Provider>
  );
};

export const useBaseContext = () => {
  return useContext(BaseContext);
};
