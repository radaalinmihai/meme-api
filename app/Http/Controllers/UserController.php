<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\User;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Auth;
use GuzzleHttp\Client;

class UserController extends Controller
{
    /**
     * login api
     *
     * @param Request $request
     * @return string[]
     * @throws \Illuminate\Validation\ValidationException
     */
    public function login(Request $request): array
    {
        $rules = [
            'username' => 'required',
            'password' => 'required',
        ];
        $this->validate($request, $rules);
        if (Auth::attempt(['username' => $request->username, 'password' => $request->password])) {
            $user = Auth::user();
            $success = $this->generateToken($user->email, $request->password);
            $success['code'] = 'OK';
            return $success;
        } else {
            return ['error' => 'Username and password don\'t match.'];
        }
    }

    private function generateToken($username, $password)
    {
        $http = new Client;
        $response = $http->post(url('/oauth/token'), [
            'allow_redirects' => false,
            'http_errors' => false,
            'form_params' => [
                'grant_type' => 'password',
                'client_id' => env('OAUTH_PASSWORD_CLIENT_ID'),
                'client_secret' => env('OAUTH_PASSWORD_CLIENT_SECRET'),
                'username' => $username,
                'password' => $password,
                'scope' => '*',
            ]
        ]);
        return json_decode((string)$response->getBody(), true);
    }

    /**
     * Register api
     *
     * @param Request $request
     * @return Response
     * @throws \Illuminate\Validation\ValidationException
     */
    public function register(Request $request): Response
    {
        $rules = [
            'username' => 'required',
            'email' => 'required|email',
            'password' => 'required',
            'c_password' => 'required|same:password',
        ];

        $this->validate($request, $rules);

        $input = $request->all();
        $input['password'] = bcrypt($input['password']);
        $user = User::create($input);
        $success = $this->generateToken($user->email, $request->password);
        $success['code'] = "OK";
        return $success;
    }

    public function refreshToken(Request $request): array
    {
        $refresh_token = $this->generateRefreshToken($request->refresh_token);
        return ['token' => $refresh_token];
    }

    private function generateRefreshToken($refreshToken)
    {
        $http = new Client;
        $response = $http->post(url('/oauth/token'), [
            'form_params' => [
                'grant_type' => 'refresh_token',
                'refresh_token' => $refreshToken,
                'client_id' => env('OAUTH_PASSWORD_CLIENT_ID'),
                'client_secret' => env('OAUTH_PASSWORD_CLIENT_SECRET'),
                'scope' => '*',
            ]
        ]);

        return json_decode((string)$response->getBody(), true);
    }

    public function details(): array
    {
        $user = Auth::user();
        return ['success' => $user];
    }

    public function logout(): array
    {
        if (Auth::check()) {
            Auth::user()->token()->revoke();
            return ['success' => 'Logged out succesfuly'];
        }
    }
}
