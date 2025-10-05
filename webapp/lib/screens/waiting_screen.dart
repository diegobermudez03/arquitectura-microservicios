import 'package:flutter/material.dart';
import '../services/sse_service.dart';
import '../models/issuer_response.dart';

class WaitingScreen extends StatefulWidget {
  const WaitingScreen({super.key});

  @override
  State<WaitingScreen> createState() => _WaitingScreenState();
}

class _WaitingScreenState extends State<WaitingScreen> with TickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _animation;
  final SseService _sseService = SseService();

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(seconds: 2),
      vsync: this,
    );
    _animation = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeInOut),
    );
    _animationController.repeat(reverse: true);
    
    // Update SSE callbacks to handle notifications in this screen
    _sseService.updateCallbacks(
      onMessage: _handleNotification,
      onError: _handleError,
    );
  }
  
  void _handleNotification(IssuerResponse response) {
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
  }
  
  void _handleError() {
  }

  @override
  void dispose() {
    _animationController.dispose();
    // Don't disconnect SSE here - it might be needed for result screens
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Processing Request'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        automaticallyImplyLeading: false,
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              AnimatedBuilder(
                animation: _animation,
                builder: (context, child) {
                  return Transform.scale(
                    scale: 0.8 + (_animation.value * 0.4),
                    child: const Icon(
                      Icons.hourglass_empty,
                      size: 120,
                      color: Colors.blue,
                    ),
                  );
                },
              ),
              const SizedBox(height: 32),
              const Text(
                'Processing your request, please wait...',
                style: TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.bold,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              const Text(
                'We are reviewing your application and will notify you of the result shortly.',
                style: TextStyle(
                  fontSize: 16,
                  color: Colors.grey,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 32),
              const CircularProgressIndicator(),
            ],
          ),
        ),
      ),
    );
  }
}
