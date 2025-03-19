function getInitAuth() {
  const init_api_token = `http://localhost:8080/auth/${process.env.API_BEARER_TOKEN}`;

  console.log(init_api_token);
  console.log(process.env.API_BEARER_TOKEN);

  return {
    type: "getInitAuth",
  };
}
export default getInitAuth;
