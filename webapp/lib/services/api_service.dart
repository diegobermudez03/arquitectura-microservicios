import 'dart:convert';
import 'package:http/http.dart' as http;
import '../env/env.dart';
import '../models/full_card.dart';

class ApiService {
  static const String _registerEndpoint = '/v1/register';
  static const String _issueEndpoint = '/v1/issue';

  static Future<Map<String, dynamic>> registerUser({
    required String firstName,
    required String lastName,
    required String countryCode,
    required String birthDate,
    required String citizenId,
  }) async {
    final url = Uri.parse('${Env.issueServiceUrl}$_registerEndpoint');
    print('Registration URL: $url');
    
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'name': firstName,
        'lastname': lastName,
        'country_code': countryCode,
        'birth_date': birthDate,
        'citizen_id': citizenId,
      }),
    );

    print('Registration API Response Status: ${response.statusCode}');
    print('Registration API Response Body: ${response.body}');
    
    if (response.statusCode == 200) {
      final decodedResponse = jsonDecode(response.body);
      print('Decoded registration response: $decodedResponse');
      return decodedResponse;
    } else {
      throw Exception('Failed to register user: ${response.statusCode} - ${response.body}');
    }
  }

  static Future<void> issueCard({
    required String userToken,
    required String cardType,
  }) async {
    final url = Uri.parse('${Env.issueServiceUrl}$_issueEndpoint');
    
    final response = await http.post(
      url,
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'user_token': userToken,
        'card_type': cardType,
      }),
    );

    if (response.statusCode == 202) {
      return;
    } else {
      throw Exception('Failed to issue card: ${response.statusCode}');
    }
  }

  static Future<List<FullCard>> getCardsByCitizenId({
    required String citizenId,
  }) async {
    final url = Uri.parse('${Env.issueServiceUrl}/v1/$citizenId/cards');
    print('Get cards URL: $url');
    
    final response = await http.get(
      url,
      headers: {'Content-Type': 'application/json'},
    );

    print('Get cards API Response Status: ${response.statusCode}');
    print('Get cards API Response Body: ${response.body}');
    
    if (response.statusCode == 200) {
      final List<dynamic> jsonList = jsonDecode(response.body);
      return jsonList.map((json) => FullCard.fromJson(json)).toList();
    } else {
      throw Exception('Failed to get cards: ${response.statusCode} - ${response.body}');
    }
  }
}
