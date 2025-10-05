import 'package:flutter/material.dart';
import 'screens/register_screen.dart';
import 'screens/card_selection_screen.dart';
import 'screens/waiting_screen.dart';
import 'screens/result_success_screen.dart';
import 'screens/result_declined_screen.dart';
import 'screens/card_search_screen.dart';
import 'models/issuer_response.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  runApp(const CardIssuanceApp());
}

class CardIssuanceApp extends StatelessWidget {
  const CardIssuanceApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Card Issuance System',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.blue),
        useMaterial3: true,
      ),
      debugShowCheckedModeBanner: false,
      initialRoute: '/',
      routes: {
        '/': (context) => const HomeScreen(),
        '/register': (context) => const RegisterScreen(),
        '/card-search': (context) => const CardSearchScreen(),
        '/card-selection': (context) {
          final userToken = ModalRoute.of(context)!.settings.arguments as String?;
          print('Card selection route - userToken: $userToken');
          if (userToken == null) {
            throw Exception('User token is null when navigating to card selection');
          }
          return CardSelectionScreen(userToken: userToken);
        },
        '/waiting': (context) => const WaitingScreen(),
        '/result-success': (context) {
          final issuedCard = ModalRoute.of(context)!.settings.arguments as IssuedCard;
          return ResultSuccessScreen(issuedCard: issuedCard);
        },
        '/result-declined': (context) {
          final declineReason = ModalRoute.of(context)!.settings.arguments as String;
          return ResultDeclinedScreen(declineReason: declineReason);
        },
      },
    );
  }
}

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Card Issuance System'),
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
      ),
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(32.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(
                Icons.credit_card,
                size: 120,
                color: Colors.blue,
              ),
              const SizedBox(height: 32),
              const Text(
                'Welcome to the Card Issuance System',
                style: TextStyle(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              const Text(
                'Apply for a credit, debit, or prepaid card in just a few simple steps.',
                style: TextStyle(
                  fontSize: 18,
                  color: Colors.grey,
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 48),
              ElevatedButton.icon(
                onPressed: () {
                  Navigator.pushNamed(context, '/register');
                },
                icon: const Icon(Icons.person_add),
                label: const Text('Start Application'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.blue,
                  foregroundColor: Colors.white,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 32,
                    vertical: 16,
                  ),
                  textStyle: const TextStyle(fontSize: 18),
                ),
              ),
              const SizedBox(height: 16),
              ElevatedButton.icon(
                onPressed: () {
                  Navigator.pushNamed(context, '/card-search');
                },
                icon: const Icon(Icons.search),
                label: const Text('Search My Cards'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.green,
                  foregroundColor: Colors.white,
                  padding: const EdgeInsets.symmetric(
                    horizontal: 32,
                    vertical: 16,
                  ),
                  textStyle: const TextStyle(fontSize: 18),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
