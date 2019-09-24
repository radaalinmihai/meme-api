<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use App\Http\Controllers\Controller;
use App\User;
use Illuminate\Support\Facades\Auth;
use GuzzleHttp\Client;

class UserController extends Controller
{
    public $successStatus = 200;
    /**
     * login api
     *
     * @return \Illuminate\Http\Response
     */
    public function login(Request $request)
    {
        $rules = [
            'username' => 'required',
            'password' => 'required',
        ];
        $this->validate($request, $rules);
        if (Auth::attempt(['username' => $request->username, 'password' => $request->password])) {
            $user = Auth::user();
            $success['token'] =  $this->generateToken($user->email, $request->password);
            return response()->json(['success' => $success], $this->successStatus);
        } else {
            return response()->json(['error' => 'Unauthorised'], 401);
        }
    }
    /**
     * Register api
     *
     * @return \Illuminate\Http\Response
     */
    public function register(Request $request)
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
        $success['token'] = $this->generateToken($user->email, $request->password);
        return response()->json(['success' => $success], $this->successStatus);
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
        return json_decode((string) $response->getBody(), true);
    }
    public function refreshToken(Request $request)
    {
        $refresh_token = $this->generateRefreshToken($request->bearerToken());
        return response()->json(['token' => $refresh_token], $this->successStatus);
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

        return json_decode((string) $response->getBody(), true);
    }
    /**
     * details api
     *
     * @return \Illuminate\Http\Response
     */
    public function checkToken()
    {
        return response()->json(['success' => true], $this->successStatus);
    }
    public function details()
    {
        $user = Auth::user();
        return response()->json(['success' => $user], $this->successStatus);
    }
    public function logout()
    {
        if(Auth::check()) {
            Auth::user()->AauthAccessToken()->delete();
            return response()->json(['success' => 'Logged out succesfuly'], $this->successStatus);
        }
    }
}
