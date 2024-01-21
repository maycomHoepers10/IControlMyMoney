import * as Yup from 'yup';

export const categoryValidationSchema = Yup.object().shape({
  name: Yup.string().required('Campo obrigatório')
});

export const transactionValidationSchema = Yup.object().shape({
  date: Yup.string().required('Campo obrigatório'),
  amount: Yup.number().required('Campo obrigatório'),
  description: Yup.string().required('Campo obrigatório')
}); 