import 'package:flutter/material.dart';
import '../services/api_service.dart';
import '../services/sse_service.dart';
import '../models/issuer_response.dart';

class CardSelectionScreen extends StatefulWidget {
  final String userToken;

  const CardSelectionScreen({super.key, required this.userToken});

  @override
  State<CardSelectionScreen> createState() => _CardSelectionScreenState();
}

class _CardSelectionScreenState extends State<CardSelectionScreen> {
  final SseService _sseService = SseService();
  bool _isLoading = false;

  Future<void> _selectCardType(String cardType) async {
    setState(() {
      _isLoading = true;
    });

    try {
      // Start listening for notifications
      _sseService.connectToNotifications(
        userToken: widget.userToken,
        onMessage: (IssuerResponse response) {
          if (mounted) {
            if (response.declineReason != null) {
              Navigator.pushReplacementNamed(
                context,
                '/result-declined',
                arguments: response.declineReason!.reason,
              );
            } else if (response.issuedCard != null) {
              Navigator.pushReplacementNamed(
                context,
                '/result-success',
                arguments: response.issuedCard,
              );
            }
          }
        },
        onError: () {
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(content: Text('Connection error. Please try again.')),
            );
            setState(() {
              _isLoading = false;
            });
          }
        },
      );
      
      await ApiService.issueCard(
        userToken: widget.userToken,
        cardType: cardType,
      );

      if (mounted) {
        Navigator.pushReplacementNamed(context, '/waiting');
      }
    } catch (e) {
      print('Error in card selection screen: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e')),
        );
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  void dispose() {
    // Don't disconnect SSE here - it needs to stay active for the waiting screen
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Select Card Type'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Choose the type of card you want to apply for:',
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 40),
            Expanded(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  _buildCardOption(
                    'Credit Card',
                    Icons.credit_card,
                    Colors.blue,
                    'credit',
                  ),
                  _buildCardOption(
                    'Debit Card',
                    Icons.account_balance_wallet,
                    Colors.green,
                    'debit',
                  ),
                  _buildCardOption(
                    'Prepaid Card',
                    Icons.payment,
                    Colors.orange,
                    'prepaid',
                  ),
                ],
              ),
            ),
            if (_isLoading)
              const Center(
                child: Column(
                  children: [
                    CircularProgressIndicator(),
                    SizedBox(height: 16),
                    Text('Processing your request...'),
                  ],
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildCardOption(String title, IconData icon, Color color, String cardType) {
    return Card(
      elevation: 4,
      child: InkWell(
        onTap: _isLoading ? null : () => _selectCardType(cardType),
        child: Container(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                icon,
                size: 64,
                color: color,
              ),
              const SizedBox(height: 16),
              Text(
                title,
                style: const TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
