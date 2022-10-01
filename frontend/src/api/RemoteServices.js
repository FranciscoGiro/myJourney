import axios from 'axios';

// Define your api url from any source.
// Pulling from your .env file when on the server or from localhost when locally
const BASE_URL = process.env.REACT_APP_API_URL;

export const login = async (username, password) => {
  try {
    const res = await axios({
      method: 'post',
      url: `${BASE_URL}/auth/login`,
      withCredentials: true,
      headers: { 'Content-Type': 'application/json' },
      data: JSON.stringify({ username, password })
    });

    console.log(res.data);

    if (res.status === 200) {
      console.log(res.data);
      return res.data;
    } else {
      return 'Something unexpected happened. Please, upload again';
    }
  } catch (error) {
    console.log(error.response.data.error);
    return error;
  }
};

export const register = async (username, email, password) => {
  try {
    const res = await axios({
      method: 'post',
      url: `${BASE_URL}/auth/register`,
      withCredentials: true,
      headers: { 'Content-Type': 'application/json' },
      data: JSON.stringify({ username, email, password })
    });

    console.log(res.data);

    if (res.status === 200) {
      return res.data;
    } else {
      return 'Something unexpected happened. Please, upload again';
    }
  } catch (error) {
    console.log(error);
    return error;
  }
};

/** No param required to retrieve all images */
/* const getAllUsers = () => {
  return axios
    .get(`${BASE_URL}/users`)
    .then((res) => {
      const data = res.results.data;
      return data;
    })
    .catch((err) => {
      console.log(err);
    });
};

/** @param {string} id
const getUser = (id) => {
  return axios
    .get(`${BASE_URL}/users/${id}`)
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      console.log(err);
    });
}; */

/** @param {string} id */
export const getImage = (id) => {
  return axios
    .get(`${BASE_URL}/images/${id}`)
    .then((res) => {
      const data = res.results.data;
      return data;
    })
    .catch((err) => {
      console.log(err);
    });
};

/** @param {file array} files */
export const uploadImages = async (files) => {
  const formData = new FormData();
  files.map(file =>
    formData.append('file', file)
  );

  try {
    const res = await axios({
      method: 'post',
      url: `${BASE_URL}/images/upload`,
      data: formData,
      headers: { 'Content-Type': 'multipart/form-data' }
    });

    if (res.status === 200) {
      return 'Images uploaded successfully';
    } else {
      return 'Something unexpected happened. Please, upload again';
    }
  } catch (error) {
    console.log(error);
    return error;
  }
};
