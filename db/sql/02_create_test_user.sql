SET @user_1 = '11111111-1111-1111-1111-111111111111';
SET @user_2 = '22222222-2222-2222-2222-222222222222';
SET @user_3 = '33333333-3333-3333-3333-333333333333';

INSERT INTO users(id, user_name, email, password)
VALUES (@user_1, 'alice', 'alice@email.com', 'alice_password')
     , (@user_2, 'bob', 'bob@email.com', 'bob_password')
     , (@user_3, 'charlie', 'charlie@email.com', 'charlie_password');
