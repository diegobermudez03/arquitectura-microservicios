import 'dart:html' as html;
import 'dart:convert';
import '../env/env.dart';
import '../models/issuer_response.dart';

class SseService {
  static final SseService _instance = SseService._internal();
  factory SseService() => _instance;
  SseService._internal();
  
  html.EventSource? _eventSource;
  String? _currentUserToken;
  Function(IssuerResponse)? _onMessageCallback;
  Function()? _onErrorCallback;
  
  void connectToNotifications({
    required String userToken,
    required Function(IssuerResponse) onMessage,
    required Function() onError,
  }) {
    // Store the callbacks
    _onMessageCallback = onMessage;
    _onErrorCallback = onError;
    _currentUserToken = userToken;
    
    // Only create a new connection if we don't have one or if the user token changed
    if (_eventSource == null || _currentUserToken != userToken) {
      disconnect(); // Close existing connection if any
      
      final url = '${Env.notificationServiceUrl}/notifications/stream?user_token=$userToken';
      _eventSource = html.EventSource(url);
    }
    
    _eventSource!.onMessage.listen((event) {
      try {
        final data = event.data;
        
        if (data == null || data.isEmpty) {
          return;
        }
        
        final jsonData = data as String;
        
        // Try to parse the data - it could be JSON or query string
        Map<String, dynamic> jsonMap;
        try {
          // First try parsing as JSON
          jsonMap = jsonDecode(jsonData) as Map<String, dynamic>;
        } catch (e) {
          // If JSON parsing fails, try as query string
          try {
            jsonMap = Map<String, dynamic>.from(
              Uri.splitQueryString(jsonData)
            );
          } catch (e2) {
            // If all parsing fails, create a simple map
            jsonMap = {'raw_data': jsonData};
          }
        }
        
        // Convert the data to proper JSON format
        final Map<String, dynamic> properJson = {};
        
        // Handle different data formats
        if (jsonMap.containsKey('raw_data')) {
          // If we have raw data, try to parse it as JSON
          try {
            final rawData = jsonMap['raw_data'] as String;
            final parsedData = Uri.splitQueryString(rawData);
            jsonMap.addAll(parsedData);
          } catch (e) {
          }
        }
        
        jsonMap.forEach((key, value) {
          if (key == 'decline_reason' && value != null) {
            // If decline_reason is already an object, use it directly
            if (value is Map<String, dynamic>) {
              properJson['decline_reason'] = value;
            } else {
              // If it's a string, wrap it in an object
              properJson['decline_reason'] = {'reason': value};
            }
          } else if (key == 'issued_card' && value != null) {
            // If issued_card is already an object, use it directly
            if (value is Map<String, dynamic>) {
              properJson['issued_card'] = value;
            } else {
              // If it's a string, try to parse it
              try {
                final cardData = Uri.splitQueryString(value.toString());
                properJson['issued_card'] = {
                  'pan': cardData['pan'] ?? '',
                  'cvv': cardData['cvv'] ?? '',
                  'expiry_date': cardData['expiry_date'] ?? '',
                  'card_type': cardData['card_type'] ?? '',
                };
              } catch (e) {
                properJson['issued_card'] = {'error': 'Failed to parse card data'};
              }
            }
          } else if (key != 'raw_data') {
            properJson[key] = value;
          }
        });
        
        
        // Check if this is just a connection status message
        if (properJson.containsKey('status') && properJson['status'] == 'connected') {
          return;
        }
        
        // Only process if we have actual card issuance data
        if (properJson.containsKey('decline_reason') || properJson.containsKey('issued_card')) {
          final issuerResponse = IssuerResponse.fromJson(properJson);
          _onMessageCallback?.call(issuerResponse);
        } else {
        }
      } catch (e) {
        _onErrorCallback?.call();
      }
    });
    
    _eventSource!.onError.listen((event) {
      _onErrorCallback?.call();
    });
    
    // Add connection state monitoring
    _eventSource!.onOpen.listen((event) {
    });
  }
  
  void updateCallbacks({
    required Function(IssuerResponse) onMessage,
    required Function() onError,
  }) {
    _onMessageCallback = onMessage;
    _onErrorCallback = onError;
    print('SSE callbacks updated');
  }
  
  bool get isConnected => _eventSource != null;
  
  void disconnect() {
    _eventSource?.close();
    _eventSource = null;
    _onMessageCallback = null;
    _onErrorCallback = null;
    _currentUserToken = null;
    print('SSE connection disconnected');
  }
}
