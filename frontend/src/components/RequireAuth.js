/* eslint-disable react/react-in-jsx-scope */
import { useLocation, Navigate, Outlet } from 'react-router-dom';
import useAuth from '../context/useAuth';

const RequireAuth = () => {
  const { auth } = useAuth();
  const location = useLocation();

  return (
    auth?.user
      ? <Outlet />
      : <Navigate to="/login" state={{ from: location }} replace />
  );
};

export default RequireAuth;
